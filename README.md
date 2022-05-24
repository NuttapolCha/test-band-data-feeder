# test-band-data-feeder

This is the first interview assignment at Band Protocol for the role 'Chan and Feeder'.

[see instruction](https://hackmd.io/@-xuuChM-TfSOB631tC1wMg/SyfnIhHLd#Simple-solution-for-this-interview)

## Solution

![sequence diagram](./docs/Simple%20Data%20Feeder%20Service.png)

### Commands

For the solution of this assignment, please run the following command at the project root path in order to 24/7 feed data from source to destination.


```sh
$./data-feeder auto-feeder
```

or

```sh
$./go run main.go auto-feeder
```

this required go 1.17 or later installed in your computer.

### Configuration

Any constant can be configured at [config.yaml](./config/config.yaml) before starting the service.

Feel free to adjust and play with it.