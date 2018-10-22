# TruChain: Backing Module

![](dep.png)

## Keeper

### Dependencies
* bank keeper
* story keeper
* category keeper

### Stores
* "backings"
    *  keys
        * `"backings:id:5"` -> `Backing`
        * `"backings:unexpired:queue"` -> `[1,2,3]`

## Notes

Every new `Backing` id is saved in a queue which is checked for expiration on each block tick.