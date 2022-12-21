## Node states
- `pStorage`: the log of operations
- `dStruct`: the data structure after applying all the operations in `pStorage` in order

## Leader states
- `lastEpoch`: epoch
- `zxid`: {epoch, counter}

## Follower states
- `lastEpochProp`: Last new epoch proposal follower f acknowledged, initially ⊥
- `lastLeaderProp`: Last new leader proposal follower f acknowledged, initially ⊥

## Transition functions
| curr_state              | incoming_message                             | next_state                                                                                          |
|-------------------------|----------------------------------------------|-----------------------------------------------------------------------------------------------------|
|                         | Broadcast{lastEpoch, Transaction{vec, zxid}} | append Transaction{vec, zxid} to `pStorage`                                                           |
| epoch == `lastLeaderProp` | Commit{epoch}                                | apply the last Transaction{vec, zxid} to `dStruct`                                                   |
| `lastEpochProp` < epoch   | NewEpoch{epoch}                              | `lastEpochProp` = epoch                                                                               |
|                         | NewLeader{epoch, hist}                       | `lastLeaderProp` = epoch<br> `pStorage` = hist<br> each transaction's epoch in `pStorage` is set to epoch |
| epoch == `lastLeaderProp` | CommitNewLeader{epoch}                       | re-apply all the Transaction{vec, zxid} in pStorage to `dStruct`                                      |
## Predicates
### node outgoing message Broadcast{`lastEpoch`, Transaction{vec, `zxid`}}
- this node is currently the leader
- node must received Store{vec} from user
- `zxid` cannot decrease

### node outgoing message Commit{content, `lastEpoch`}
- this node is currently the leader
- `lastEpoch` cannot decrease
- node must received valid AckTxn{} reply from a quorum of followers