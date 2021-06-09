package noolite

import (
	"errors"
	"fmt"
)

type Response []byte

func NewResponse(data []byte) (Response, error) {
	if !isCorrect(data) {
		return nil, errors.New("wrong data")
	}

	return data, nil
}

func isCorrect(data []byte) bool {
	if len(data) < respLen {
		return false
	}
	if data[0] != RespSt {
		return false
	}
	if data[16] != RespSp {
		return false
	}

	return calcCrc(data) == data[15]
}

func (r Response) String() string {
	return fmt.Sprintf(
		"ST: %d, MODE: %d, CTR: %d, TOGL: %d, CH: %d, CMD: %d, FMT: %d, "+
			"D0: %d, D1: %d, D2: %d, D3: %d, ID0: %d, ID1: %d, ID2: %d, ID3: %d, CRC: %d, SP: %d",
		int(r[0]),
		int(r[1]),
		int(r[2]),
		int(r[3]),
		int(r[4]),
		int(r[5]),
		int(r[6]),
		int(r[7]),
		int(r[8]),
		int(r[9]),
		int(r[10]),
		int(r[11]),
		int(r[12]),
		int(r[13]),
		int(r[14]),
		int(r[15]),
		int(r[16]),
	)
}

func (r Response) GetChannel() int {
	return int(r[4])
}

func (r Response) GetCommand() Cmd {
	return Cmd(r[5])
}

func (r Response) GetD2() int {
	return int(r[9])
}

func (r Response) IsSuccess() bool {
	ctr := RespCtr(r[2])

	return ctr == RespCtrDone || ctr == RespCtrBinded
}
