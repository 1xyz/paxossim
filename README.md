# paxos-simulation

This is a simple simulation attempting to mimic the paxos protocol under varying conditions

## Replica

**Invariants**

- R1: There are no two commands decided for the same slot. 
    - i.e given two Replicas r1 & r2, two commands c1 and c2. for a 
given slot s, if r1.Decisions contains (s, c1) and r2.Decisions contains (s, c2), then c1 and c2 are the same command
- R2: All commands, upto slot_out are the set of decisions.
- R3: For all replicas r, r.state is essentially the result of applying the commands in the set of decisions 
(R.Decisions) in order
- R4: For each replica, r.slot_out cannot decrease over time
- R5: A replica proposes commands for configuration it knows about. 
    - when the slot s is in *[slot_in, slot_out+WINDOW)*
    
 ## Acceptor
 
 **Invariants**
 
 - A1: An Acceptor can only adopt strictly increasing ballot numbers
 - A2: An Acceptor a can only adopt a p-value: <b, s, c> if its currently adopted ballot_number b is the
 same as that of the p-value.
    - i.e p_value.b = a.b, for  <b, s, c> to be accepted
 - A3: An Acceptor a cannot remove values from its accepted list. 
    - This is an impractical invariant, since in every phase1 response includes the entire Accepted list to a Scout
 - A4: For any two acceptors a and a', For the same ballot_number, slot_number combination accepted, there can only be 
 one proposed command associated
    - i.e if a.accepted contains <b, s, c> & a'.accepted contains <b, s, c'> then. c = c'
 - A5: If a majority of acceptors have accepted a p-value for the current ballot, slot and command, then any future ballots
 proposing can only propose the same command and slot combination
    - i.e for a majority of acceptors a. if <b, s, c> are in a.accepted. If a new ballot <b', s, c'> has been accepted 
    for a' then c = c'.