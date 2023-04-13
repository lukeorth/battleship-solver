package battleshipsolver

import (
	"strconv"
)

const (
    CARRIER_NAME = "carrier"
    BATTLESHIP_NAME = "battleship"
    SUBMARINE_NAME = "submarine"
    CRUISER_NAME = "cruiser"
    DESTROYER_NAME = "destroyer"

    CARRIER_MASK = "1111100000000000"
    BATTLESHIP_MASK = "1111000000000000"
    SUBMARINE_MASK = "1110000000000000"
    CRUISER_MASK = "1110000000000000"
    DESTROYER_MASK = "1100000000000000"

    CARRIER_LENGTH = 5
    BATTLESHIP_LENGTH = 4
    SUBMARINE_LENGTH = 3
    CRUISER_LENGTH = 3
    DESTROYER_LENGTH = 2
)

type Ship struct {
    Name string
    Mask uint
    Length int
}

type Fleet struct {
    ships map[string]*Ship
}

func BuildFleet() *Fleet {
    carrierMask, _ := strconv.ParseUint(CARRIER_MASK, 2, 64)
    battleshipMask, _ := strconv.ParseUint(BATTLESHIP_MASK, 2, 64)
    submarineMask, _ := strconv.ParseUint(SUBMARINE_MASK, 2, 64)
    cruiserMask, _ := strconv.ParseUint(CRUISER_MASK, 2, 64)
    destroyerMask, _ := strconv.ParseUint(DESTROYER_MASK, 2, 64)

    ships := map[string]*Ship{
        CARRIER_NAME: {CARRIER_NAME, uint(carrierMask), CARRIER_LENGTH},
        BATTLESHIP_NAME: {BATTLESHIP_NAME, uint(battleshipMask), BATTLESHIP_LENGTH},
        SUBMARINE_NAME: {SUBMARINE_NAME, uint(submarineMask), SUBMARINE_LENGTH},
        CRUISER_NAME: {CRUISER_NAME, uint(cruiserMask), CRUISER_LENGTH},
        DESTROYER_NAME: {DESTROYER_NAME, uint(destroyerMask), DESTROYER_LENGTH},
    }

    return &Fleet{ships: ships}
}

func (f *Fleet) GetBiggestShipSize() int {
    biggest := 0
    for _, s := range f.ships {
        if s.Length > biggest {
            biggest = s.Length
        }
    }

    return biggest
}

func (f *Fleet) SinkShip(ship string) {
    delete(f.ships, ship)
}
