package fieldtrie

import (
	"reflect"
	"sync"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/beacon-chain/state/stateutil"
	"github.com/prysmaticlabs/prysm/beacon-chain/state/types"
)

// FieldTrie is the representation of the representative
// trie of the particular field.
type FieldTrie struct {
	*sync.RWMutex
	reference   *stateutil.Reference
	fieldLayers [][]*[32]byte
	field       types.FieldIndex
	dataType    types.DataType
	length      uint64
	numOfElems  int
}

// NewFieldTrie is the constructor for the field trie data structure. It creates the corresponding
// trie according to the given parameters. Depending on whether the field is a basic/composite array
// which is either fixed/variable length, it will appropriately determine the trie.
func NewFieldTrie(field types.FieldIndex, dataType types.DataType, elements interface{}, length uint64) (*FieldTrie, error) {
	if elements == nil {
		return &FieldTrie{
			field:      field,
			dataType:   dataType,
			reference:  stateutil.NewRef(1),
			RWMutex:    new(sync.RWMutex),
			length:     length,
			numOfElems: 0,
		}, nil
	}
	fieldRoots, err := fieldConverters(field, []uint64{}, elements, true)
	if err != nil {
		return nil, err
	}
	if err := validateElements(field, dataType, elements, length); err != nil {
		return nil, err
	}
	switch dataType {
	case types.BasicArray:
		fl, err := stateutil.ReturnTrieLayer(fieldRoots, length)
		if err != nil {
			return nil, err
		}
		return &FieldTrie{
			fieldLayers: fl,
			field:       field,
			dataType:    dataType,
			reference:   stateutil.NewRef(1),
			RWMutex:     new(sync.RWMutex),
			length:      length,
			numOfElems:  reflect.Indirect(reflect.ValueOf(elements)).Len(),
		}, nil
	case types.CompositeArray, types.CompressedArray:
		return &FieldTrie{
			fieldLayers: stateutil.ReturnTrieLayerVariable(fieldRoots, length),
			field:       field,
			dataType:    dataType,
			reference:   stateutil.NewRef(1),
			RWMutex:     new(sync.RWMutex),
			length:      length,
			numOfElems:  reflect.Indirect(reflect.ValueOf(elements)).Len(),
		}, nil
	default:
		return nil, errors.Errorf("unrecognized data type in field map: %v", reflect.TypeOf(dataType).Name())
	}

}

// RecomputeTrie rebuilds the affected branches in the trie according to the provided
// changed indices and elements. This recomputes the trie according to the particular
// field the trie is based on.
func (f *FieldTrie) RecomputeTrie(indices []uint64, elements interface{}) ([32]byte, error) {
	f.Lock()
	defer f.Unlock()
	var fieldRoot [32]byte
	if len(indices) == 0 {
		return f.TrieRoot()
	}
	fieldRoots, err := fieldConverters(f.field, indices, elements, false)
	if err != nil {
		return [32]byte{}, err
	}
	if err := f.validateIndices(indices); err != nil {
		return [32]byte{}, err
	}
	switch f.dataType {
	case types.BasicArray:
		fieldRoot, f.fieldLayers, err = stateutil.RecomputeFromLayer(fieldRoots, indices, f.fieldLayers)
		if err != nil {
			return [32]byte{}, err
		}
		f.numOfElems = reflect.ValueOf(elements).Len()
		return fieldRoot, nil
	case types.CompositeArray:
		fieldRoot, f.fieldLayers, err = stateutil.RecomputeFromLayerVariable(fieldRoots, indices, f.fieldLayers)
		if err != nil {
			return [32]byte{}, err
		}
		f.numOfElems = reflect.ValueOf(elements).Len()
		return stateutil.AddInMixin(fieldRoot, uint64(len(f.fieldLayers[0])))
	case types.CompressedArray:
		numOfElems, err := f.field.ElemsInChunk()
		if err != nil {
			return [32]byte{}, err
		}
		// We remove the duplicates here in order to prevent
		// duplicated insertions into the trie.
		newIndices := []uint64{}
		indexExists := make(map[uint64]bool)
		newRoots := make([][32]byte, 0, len(fieldRoots)/int(numOfElems))
		for i, idx := range indices {
			startIdx := idx / numOfElems
			if indexExists[startIdx] {
				continue
			}
			newIndices = append(newIndices, startIdx)
			indexExists[startIdx] = true
			newRoots = append(newRoots, fieldRoots[i])
		}
		fieldRoot, f.fieldLayers, err = stateutil.RecomputeFromLayerVariable(newRoots, newIndices, f.fieldLayers)
		if err != nil {
			return [32]byte{}, err
		}
		f.numOfElems = reflect.ValueOf(elements).Len()
		return stateutil.AddInMixin(fieldRoot, uint64(f.numOfElems))
	default:
		return [32]byte{}, errors.Errorf("unrecognized data type in field map: %v", reflect.TypeOf(f.dataType).Name())
	}
}

// CopyTrie copies the references to the elements the trie
// is built on.
func (f *FieldTrie) CopyTrie() *FieldTrie {
	if f.fieldLayers == nil {
		return &FieldTrie{
			field:      f.field,
			dataType:   f.dataType,
			reference:  stateutil.NewRef(1),
			RWMutex:    new(sync.RWMutex),
			length:     f.length,
			numOfElems: f.numOfElems,
		}
	}
	dstFieldTrie := make([][]*[32]byte, len(f.fieldLayers))
	for i, layer := range f.fieldLayers {
		dstFieldTrie[i] = make([]*[32]byte, len(layer))
		copy(dstFieldTrie[i], layer)
	}
	return &FieldTrie{
		fieldLayers: dstFieldTrie,
		field:       f.field,
		dataType:    f.dataType,
		reference:   stateutil.NewRef(1),
		RWMutex:     new(sync.RWMutex),
		length:      f.length,
		numOfElems:  f.numOfElems,
	}
}

// TrieRoot returns the corresponding root of the trie.
func (f *FieldTrie) TrieRoot() ([32]byte, error) {
	switch f.dataType {
	case types.BasicArray:
		return *f.fieldLayers[len(f.fieldLayers)-1][0], nil
	case types.CompositeArray:
		trieRoot := *f.fieldLayers[len(f.fieldLayers)-1][0]
		return stateutil.AddInMixin(trieRoot, uint64(len(f.fieldLayers[0])))
	case types.CompressedArray:
		trieRoot := *f.fieldLayers[len(f.fieldLayers)-1][0]
		return stateutil.AddInMixin(trieRoot, uint64(f.numOfElems))
	default:
		return [32]byte{}, errors.Errorf("unrecognized data type in field map: %v", reflect.TypeOf(f.dataType).Name())
	}
}

// FieldReference returns the underlying field reference
// object for the trie.
func (f *FieldTrie) FieldReference() *stateutil.Reference {
	return f.reference
}

// Empty checks whether the underlying field trie is
// empty or not.
func (f *FieldTrie) Empty() bool {
	return f == nil || len(f.fieldLayers) == 0
}

// InsertFieldLayer manually inserts a field layer. This method
// bypasses the normal method of field computation, it is only
// meant to be used in tests.
func (f *FieldTrie) InsertFieldLayer(layer [][]*[32]byte) {
	f.fieldLayers = layer
}
