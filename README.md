# Near-Cities
Get all surrounding cities near the target city within
the specified radius (miles).

## Build
Execute the specific command mentioned in ```builds.txt``` depending
on the target environment.

## Usage
There are 2 ways to execute this app using CLI:

### Direct Call with Args

```
main.exe [city name] [radius] [save to a JSON file(y/n)]
```

### Call (args are captured in friendly way)
```
main.exe
```

## Result
The result is a list of objects with the following properties:

*   ```city``` : City Name (string)

*   ```state``` : State Name (string)

*   ```lat``` : Latitude (float64)

*   ```lng``` : Longitude (float64)

*   ```density``` : Population Density of The City (Int)

*   ```timezone``` : The Timezone of The City (string)

## Copyright and license
Copyright 2022 Yaseen Al Mufti. Code released under [MIT License](https://github.com/YaseenAlMufti/near-cities/blob/main/LICENSE).
