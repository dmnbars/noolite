package noolite

import "fmt"

type Command []byte

func NewCommand(
	mode Mode,
	ctr CommandCtr,
	channel int,
	cmd Cmd,
) Command {
	data := Command{
		CommandSt,
		byte(mode),
		byte(ctr),
		0,
		byte(channel),
		byte(cmd),
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		CommandSp,
	}

	data[15] = calcCrc(data)

	return data
}

func (c Command) String() string {
	return fmt.Sprintf(
		"ST: %d, MODE: %d, CTR: %d, RES: %d, CH: %d, CMD: %d, FMT: %d, "+
			"D0: %d, D1: %d, D2: %d, D3: %d, ID0: %d, ID1: %d, ID2: %d, ID3: %d, CRC: %d, SP: %d",
		int(c[0]),
		int(c[1]),
		int(c[2]),
		int(c[3]),
		int(c[4]),
		int(c[5]),
		int(c[6]),
		int(c[7]),
		int(c[8]),
		int(c[9]),
		int(c[10]),
		int(c[11]),
		int(c[12]),
		int(c[13]),
		int(c[14]),
		int(c[15]),
		int(c[16]),
	)
}
