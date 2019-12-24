# paxos-simulation

This is a simple simulation attempting to mimic the paxos protocol under varying conditions

## Replica

**Invariants**

- R1 - There are no two command decided for the same slot. 
    - i.e given two Replicas r1 & r2, two commands c1 and c2. for a 
given slot s, if r1.Decisions contains (s, c1) and r2.Decisions contains (s, c2), then c1 and c2 are the same command
- R2 - All commands, upto slot_out are the set of decisions.
- R3 - For all replicas r, r.state is essentially the result of applying the commands in the set of decisions 
(R.Decisions) in order
- R4 - For each replica, r.slot_out cannot decrease over time
- R5 - A replica proposes commands for configuration it knows about. 
    - when the slot s is in *[slot_in, slot_out+WINDOW)*