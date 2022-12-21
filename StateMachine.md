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
### Follower
| curr_state              | incoming_message                             | next_state                                                                                          |
|-------------------------|----------------------------------------------|-----------------------------------------------------------------------------------------------------|
| `lastLeaderProp` == epoch   | <Broadcast, PropTxn{epoch, Txn{vec, zxid}}> | append PropTxn{epoch, Txn{vec, zxid}} to `pStorage`                                                           |
| `lastLeaderProp` == epoch | <Commit, CommitTxn{epoch}>                                | apply the latest transaction in `pStorage` to `dStruct`                                                   |

### Candidate
| curr_state              | incoming_message                             | next_state                                                                                          |
|-------------------------|----------------------------------------------|-----------------------------------------------------------------------------------------------------|
| `lastEpochProp` < epoch   | <NewEpoch, Epoch{epoch}>                             | `lastEpochProp` = epoch                                                                               |
| `lastEpochProp` == epoch   | <NewLeader, EpochHist{epoch, PropTxn[]}>                      | `lastLeaderProp` = epoch<br>`pStorage` = PropTxn[]<br>each transaction's epoch in `pStorage` is set to epoch |
| `lastLeaderProp` == epoch | <CommitNewLeader, Epoch{epoch}>                      | re-apply all the transactions in `pStorage` to `dStruct`                                      |
## Predicates
### node outgoing message <Broadcast, PropTxn{`lastEpoch`, Txn{vec, `zxid`}}>
- node must received Store{vec} from user
- for all prior <Broadcast, PropTxn{epoch, Txn{vec, zxid}}> AND <Commit, CommitTxn{content, epoch}> sent, `lastEpoch` >= epoch AND `zxid` >= zxid
- inference: this node is the current leader, with zxid `zxid`

### node outgoing message <Commit, CommitTxn{content, `lastEpoch`}>
- node must received valid AckTxn{} reply from a quorum of followers
- for all prior <Broadcast, PropTxn{epoch, Txn{vec, zxid}}> AND <Commit, CommitTxn{content, epoch}> sent, `lastEpoch` >= epoch
- inference: this node is the current leader, with epoch `lastEpoch`

### node outgoing message <NewLeader, EpochHist{`lastEpoch`, PropTxn[]}>
- for all {EpochHist{epoch', PropTxn'[]} reply the node received | `lastEpoch` == epoch'}, PropTxn[] is as new as or newer than PropTxn'[]

### node outgoing message <CommitNewLeader, Epoch{`lastEpoch`}>
- for all prior <CommitNewLeader, Epoch{epoch}> sent, `lastEpoch` > epoch
- inference: this node will be the next leader
- assertion: no two leaders can exist with the same epoch