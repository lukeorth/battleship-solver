package battleshipsolver

import "fmt"

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
    sunkAt Locator
}

type fleet struct {
    ships map[string]*ship
    hitCount int
}

func buildFleet() *fleet {
    ships := map[string]*ship{
        Carrier: {Carrier, carrierMask, carrierLength, nil},
        Battleship: {Battleship, battleshipMask, battleshipLength, nil},
        Submarine: {Submarine, submarineMask, submarineLength, nil},
        Cruiser: {Cruiser, cruiserMask, cruiserLength, nil},
        Destroyer: {Destroyer, destroyerMask, destroyerLength, nil},
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

func (f *fleet) sink(position Locator, shipName string) {
    fmt.Println("INSIDE OF SHIP SINK")
    ship := f.ships[shipName]
    fmt.Println("GOT SHIP")
    f.hitCount -= ship.length
    fmt.Println("DECREASED HIT COUNT")
    ship.sunkAt = position
    fmt.Println("SET SHIP SUNK AT")
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
        if ship.sunkAt == nil {
            ships[ship.name] = ship
        }
    }
    return ships
}

func (f *fleet) sunkShips() []*ship {
    ships := make([]*ship, 0, 5)
    for _, ship := range f.ships {
        if ship.sunkAt != nil {
            ships = append(ships, ship)
        }
    }
    return ships
}
