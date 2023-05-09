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
    sunk bool
}

type fleet struct {
    ships map[string]*ship
    hitCount int
}

func buildFleet() *fleet {
    ships := map[string]*ship{
        Carrier: {Carrier, carrierMask, carrierLength, false},
        Battleship: {Battleship, battleshipMask, battleshipLength, false},
        Submarine: {Submarine, submarineMask, submarineLength, false},
        Cruiser: {Cruiser, cruiserMask, cruiserLength, false},
        Destroyer: {Destroyer, destroyerMask, destroyerLength, false},
    }
    return &fleet{
        ships: ships,
        hitCount: 0,
    }
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

func (f *fleet) sink(shipName string) {
    ship := f.ships[shipName]
    f.hitCount -= ship.length
    ship.sunk = true
}

func (f *fleet) remove(shipName string) {
    delete(f.ships, shipName)
}

func (f *fleet) hit() {
    f.hitCount += 1
}

func (f *fleet) floatingShips() map[string]*ship {
    ships := make(map[string]*ship)
    for _, ship := range f.ships {
        if !ship.sunk {
            ships[ship.name] = ship
        }
    }
    return ships
}

func (f *fleet) sunkShips() map[string]*ship {
    ships := make(map[string]*ship)
    for _, ship := range f.ships {
        if ship.sunk {
            ships[ship.name] = ship
        }
    }
    return ships
}
