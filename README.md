# eaglesight-backend

This is the back-end of EagleSight.

Here is the [Trello board](https://trello.com/b/FcGCRZGN/eaglesight)

## How to run it:

You will need two files next to the eaglesight-backend's binary: 

- `map.esmap`: The file containing the map. This is required for the collisions' calculations. In developement, this is usually a symlink to an `.esmap` file in the [props](https://github.com/EagleSight/EagleSight-props/tree/master/map)

- `players.json`: list of all the players _registered_ for the on this server game.
