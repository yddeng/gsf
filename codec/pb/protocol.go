package pb

import (
	"fmt"
	"reflect"
)

type Serializer interface {
	//反序列化
	Unmarshal(data []byte, o interface{}) (err error)
	//序列化
	Marshal(data interface{}) ([]byte, error)
}

type Protocol struct {
	id2Type map[uint16]reflect.Type
	type2Id map[reflect.Type]uint16
	serial  Serializer
}

func NewProtocol(serial Serializer) *Protocol {
	return &Protocol{
		id2Type: map[uint16]reflect.Type{},
		type2Id: map[reflect.Type]uint16{},
		serial:  serial,
	}
}

func (this *Protocol) Register(id uint16, msg interface{}) {
	tt := reflect.TypeOf(msg)

	if _, ok := this.id2Type[id]; ok {
		panic(fmt.Sprintf("%d already register to type:%s\n", id, tt))
	}

	this.id2Type[id] = tt
	this.type2Id[tt] = id
}

func (this *Protocol) Marshal(data interface{}) (uint16, []byte, error) {
	id, ok := this.type2Id[reflect.TypeOf(data)]
	if !ok {
		return 0, nil, fmt.Errorf("codec.pb:Marshal type: %s undefined", reflect.TypeOf(data))
	}

	ret, err := this.serial.Marshal(data)
	if err != nil {
		return 0, nil, err
	}

	return id, ret, nil
}

func (this *Protocol) Unmarshal(msgID uint16, data []byte) (msg interface{}, err error) {
	tt, ok := this.id2Type[msgID]
	if !ok {
		err = fmt.Errorf("codec.pb:Unmarshal msgID: %d undefined", msgID)
		return
	}

	//反序列化的结构
	msg = reflect.New(tt.Elem()).Interface()
	err = this.serial.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
