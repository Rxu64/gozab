# TODOs
- update upFollowers and the cleanup of  election-messenger routine, also need to redesign the heartbeat-to-dead logics
- do the global cleanup in upFollowersUpdateRoutine quorum dead
- leader needs a way to know quorum dead and thus enter election
- 

# Building Blocks

## Go Routines
**Definition**: functions called through go and have internal non-stop loops

**Caution**: need to manually cleanup (return) all routines

**TODO**: a cross means the cleanup procedures are not implemented yet  

- Election
    - upFollowersUpdateRoutine ( x )
    - ElectionMessengerRoutine ( x )

- Leader
    - MessengerRoutine ( x ) subject to heartbeat 
        - BeatSender ( x )
        - PropAndCmtRoutine ( x )
    - AckToCmtRoutine ( x )

## Concurrent Global Variables
**Caution**: need to update through a dedicated routine, which provide a chanel to other parts

_Constants excluded_
- upFollowers
- upNum

