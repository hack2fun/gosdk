# go-sdk

## demo
[see demo](./demo/main.go)

## function list

- [Account](./account.go)
  - GetAccountByAddr
  - BroadcastTx
- [Bank](./bank.go)
  - GetBalance
  - GetBalanceList
  - Send
  - MultiSend
  - MultiSendWithDiffAmount
- [Validator](./validator.go)
  - GetValidator
  - GetValidatorList
- [Delegate](./delegate.go)
  - QueryDelegatorDelegations
  - QueryDelegateReward
  - WithdrawDelegatorReward
  - DelegateVeToken
  - DelegateCGT
  - UnDelegateCGT
- [Exchange](./exchange.go)
  - ExchangeToGovToken
  - ExchangeToPlatformToken
- [utils](./utils.go)
  - ConvertAddress
  - ConvertToCysicAddress
  - ConvertToETHAddress
