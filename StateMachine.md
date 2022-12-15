## Node states
- `dStruct`
- `pStorage`

## Leader (for convenient only) states
- `lastEpoch`
- `zxid`

## Follower states
- `lastEpochProp`
- `lastLeaderProp`

## Transition functions
| curr_state              | incoming_message                             | next_state                                                                                          |
|-------------------------|----------------------------------------------|-----------------------------------------------------------------------------------------------------|
|                         | NewLeader{epoch, hist}                       | `lastLeaderProp` = epoch<br> `pStorage` = hist<br> each transaction's epoch in `pStorage` is set to epoch |
|                         | Broadcast{lastEpoch, Transaction{vec, zxid}} | append Transaction{vec, zxid} to `pStorage`                                                           |
| epoch == `lastLeaderProp` | Commit{epoch}                                | apply the last Transaction{vec, zxid} to `dStruct`                                                   |
| epoch == `lastLeaderProp` | CommitNewLeader{epoch}                       | re-apply all the Transaction{vec, zxid} in pStorage to `dStruct`                                      |
| `lastEpochProp` < epoch   | NewEpoch{epoch}                              | `lastEpochProp` = epoch                                                                               |
## Predicates
### leader outgoing message Broadcast{`lastEpoch`, Transaction{vec, `zxid`}} to target follower
- leader must received Store{vec} from user
- the outgoing lastEpoch should match with `lastEpoch`
- `zxid` cannot decrease

### leader outgoing message Commit{content, `lastEpoch`} to target follower
- the outgoing lastEpoch should match with `lastEpoch`