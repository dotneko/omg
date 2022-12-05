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
  alias_length: 3

chain:
  daemon: 'onomyd'
  chain_id: 'onomy-testnet-1'
  address_prefix: 'onomy'
  valoper_prefix: 'onomyvaloper'
  base_denom: 'anom'
  token: 'nom'
  decimals: 18

options:
  default_fee: '1212anom'
  gas_adjust: 1.2
  keyring_backend: 'test'
  remainder: '1nom'
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

Import from keyring pass

```
omg a i --keyring pass                  # omg a i -k pass
```

2. Show addresses

```
omg address show                        # omg a s
```

3. Show active validators

```
omg validator show                      # omg v s
```

4. Get valoper address for validator moniker

```
omg validator show [moniker]            # omg v s nomblocks
```

5. Check balances (all accounts in address book)

```
omg balances --all                      # omg b -a
```

6. Check rewards (all accounts in address book)

```
omg rewards --all                       # omg r -a
```

7. Withdraw all rewards

```
omg tx withdraw-rewards [delegator]     # omg tx wd user1
```

8. Delegate 10nom to validator

```
omg tx delegate [delegator] [moniker] 10nom
```

9. Delegate all to validator, leaving 1000000000000anom minimum remainder

```
omg tx delegate [delegator] [moniker] --remainder 1000000000000anom
```

10. Restake all leaving 10nom (withdraw all then delegate)
```
omg tx restake [delegator] [moniker] --remainder 10nom
```

11. Automatic confirm without prompt using `--yes` or `-y` flag

```
omg tx restake [delegator] [moniker] -r 10nom -y
```

12. Check validator commissions

```
omg validator commissions [moniker]               # omg v c nomblocks
```

13. Withdraw commissions and rewards for validator

```
omg tx withdraw-commission [name] [moniker]       # omg tx wc validator nomblocks
```

14. Restake commissions and self-delegate rewards for validator

```
omg tx restake [validator] [validator moniker|valoper-address] --commission -r 100nom
```