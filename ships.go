package battleshipsolver

const (
    Carrier = "carrier"
    Battleship = "battleship"
    Submarine = "submarine"
    Cruiser = "cruiser"
    Destroyer = "destroyer"

    carrierMask uint = 63488       // 1111100000000000
    battleshipMask uint = 61440    //1111000000000000
    submarineMask uint = 57344     // 1110000000000000
    cruiserMask uint = 57344       // 1110000000000000
    destroyerMask uint = 49152     // 1100000000000000
    
    carrierLength = 5
    battleshipLength = 4
    submarineLength = 3
    cruiserLength = 3
    destroyerLength = 2
)

type ship struct {
    name string
    mask uint
    length int
}

type fleet struct {
    ships map[string]*ship
}

func buildFleet() *fleet {
    ships := map[string]*ship{
        Carrier: {Carrier, carrierMask, carrierLength},
        Battleship: {Battleship, battleshipMask, battleshipLength},
        Submarine: {Submarine, submarineMask, submarineLength},
        Cruiser: {Cruiser, cruiserMask, cruiserLength},
        Destroyer: {Destroyer, destroyerMask, destroyerLength},
    }

    return &fleet{ships: ships}
}

func (f *fleet) getBiggestShipSize() int {
    biggest := 0
    for _, s := range f.ships {
        if s.length > biggest {
            biggest = s.length
        }
    }

    return biggest
}

func (f *fleet) sinkShip(ship string) {
    delete(f.ships, ship)
}
