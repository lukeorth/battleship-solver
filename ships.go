package battleshipsolver

const (
    Carrier = "carrier"
    Battleship = "battleship"
    Submarine = "submarine"
    Cruiser = "cruiser"
    Destroyer = "destroyer"

    carrierMask uint = 992       // 1111100000
    battleshipMask uint = 960    // 1111000000
    submarineMask uint = 896     // 1110000000
    cruiserMask uint = 896       // 1110000000
    destroyerMask uint = 768     // 1100000000
    
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
