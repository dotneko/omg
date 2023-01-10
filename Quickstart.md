Quickstart Guide
================

## Installation from binary

Download latest release from https://github.com/dotneko/omg/releases

Archive contains:
* `omg` binary
* `.omgconfig.yaml` configuration file

After extracting, ensure `omg` in path and move/copy `.omgconfig.yaml` to $HOME directory

## Configuration in `.omgconfig.yaml`

```
app:
  addrbook_path: '$HOME'
  addrbook_filename: '.omg.json'
  min_alias_length: 3

chain:
  daemon: 'onomyd'
  daemon_path: '/path/to/bin'
  #daemon_path: '$CAN_BE_ENVIRONMENT_VARIABLE'
  chain_id: 'onomy-testnet-1'
  address_prefix: 'onomy'
  valoper_prefix: 'onomyvaloper'
  base_denom: 'anom'
  token: 'NOM'
  decimals: 18

options:
  default_fee: '1000anom'
  gas_adjust: 1.2
  keyring_backend: 'test'
  remainder: '1NOM'
```

Important parameters:
* `chain_id`: configure to match current chain
* `default_fee`: fee configuration when executing transactions
* `keyring_backend`: set to keyring-store, e.g. 'pass' as alternative
* `remainder`: amount in *nom* or *anom* to be minimum left after restaking/delegating *all*

## Basic Usage

1. Import name/address from keyring

```
omg address import                      # omg a i
```

2. Import from keyring pass

```
omg a i --keyring pass                  # omg a i -k pass
```

3. Show addresses

```
omg address show                        # omg a s
```

4. Show active validators

```
omg validator show                      # omg v s
```

5. Get valoper address for validator moniker

```
omg validator show [moniker]            # omg v s nomblocks
```

6. Check balances (all accounts in address book)

```
omg balances --all                      # omg b -a
```

7. Check rewards (all accounts in address book)

```
omg rewards --all                       # omg r -a
```

8. Withdraw all rewards

```
omg tx withdraw-rewards [delegator]     # omg tx wd user1
```

9. Delegate 10nom to validator

```
omg tx delegate [delegator] [moniker] 10nom
```

10. Delegate all to validator, leaving 1000000000000anom minimum remainder

```
omg tx delegate [delegator] [moniker] --remainder 1000000000000anom
```

11. Restake all leaving 10nom (withdraw all then delegate)
```
omg tx restake [delegator] [moniker] --remainder 10nom
```

12. Automatic confirm without prompt using `--yes` or `-y` flag

```
omg tx restake [delegator] [moniker] -r 10nom -y
```

13. Check validator commissions

```
omg validator commissions [moniker]               # omg v c nomblocks
```

14. Withdraw commissions and rewards for validator

```
omg tx withdraw-commission [name] [moniker]       # omg tx wc validator nomblocks
```

15. Restake commissions and self-delegate rewards for validator

```
omg tx restake [validator] [validator moniker|valoper-address] --commission -r 100nom
```

16. Send from user [name] to address

```
omg tx send [name] onomy146vvft2t99hdylzqsfccuugfw3eplar7vu9t8a
```

17. Conversion tool

```
omg c 1.5nom          # Returns 1500000000000000000anom
```

```
omg c 500_000anom     # Returns 0.0000000000005nom
```

18. Query delegator bonded amount and shares

```
omg delegation [name] [moniker]                   # omg dlg validator nomblocks
```