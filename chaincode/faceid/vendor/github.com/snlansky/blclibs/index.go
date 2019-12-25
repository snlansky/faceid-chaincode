package blclibs

import (
	"fmt"
	"strconv"
)

type Index struct {
	app   string
	table string
	key   string
	ck    *Table
}

func NewIndex(app, table, key string) *Index {
	ck := NewTable(fmt.Sprintf("%s_%s_%s_Index", app, table, key), "address", "index")
	return &Index{app: app, table: table, key: key, ck: ck}
}

func (index *Index) Save(stub IContractStub, addr Address, value []byte) (int, error) {
	count, err := index.Total(stub, addr)
	if err != nil {
		return 0, err
	}

	idx := strconv.Itoa(count)
	currentCount := strconv.Itoa(count + 1)

	// Address_N : TicketID
	err = index.ck.Insert(stub, []string{string(addr), idx}, value)
	if err != nil {
		return 0, err
	}

	// Address: Count
	err = stub.PutState(index.makeCountKey(addr), []byte(currentCount))
	return count + 1, err
}

func (index *Index) Total(stub IContractStub, address Address) (int, error) {
	key := index.makeCountKey(address)
	countBytes, err := stub.GetState(key)
	if err != nil {
		return 0, err
	}

	if countBytes == nil || len(countBytes) == 0 {
		return 0, nil
	}

	return strconv.Atoi(string(countBytes))
}

func (index *Index) List(stub IContractStub, address Address, offset, limit int, order bool) ([][]byte, error) {
	count, err := index.Total(stub, address)
	if err != nil {
		return nil, err
	}

	var list [][]byte

	if count <= offset {
		return list, nil
	}

	j := 1
	if order {
		for i := 0; i < count; i++ {
			if limit > 0 && j > limit {
				break
			}
			value, err := index.getValue(stub, address, i)
			if err != nil {
				return nil, err
			}
			list = append(list, value)
			j++
		}
	} else {
		for i := count - offset - 1; i >= 0; i-- {
			if limit > 0 && j > limit {
				break
			}
			value, err := index.getValue(stub, address, i)
			if err != nil {
				return nil, err
			}
			list = append(list, value)
			j++
		}
	}

	return list, nil
}

func (index *Index) Filter(stub IContractStub, address Address, f func(value []byte) (bool, error)) error {
	count, err := index.Total(stub, address)
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		value, err := index.getValue(stub, address, i)
		if err != nil {
			return err
		}
		ctiu, err := f(value)
		if err != nil {
			return err
		}
		if !ctiu {
			return nil
		}
	}

	return nil
}

func (index *Index) getValue(stub IContractStub, address Address, idx int) ([]byte, error) {
	return index.ck.GetValue(stub, []string{string(address), strconv.Itoa(idx)})
}

func (index *Index) makeCountKey(address Address) string {
	return fmt.Sprintf("%s_%s_%s_%s_Total", address, index.app, index.table, index.key)
}
