# paxos-simulation

![alt text](https://upload.wikimedia.org/wikipedia/commons/e/e9/Gaios_panoramic.jpg  "Paxos, Greece")

This is a simple simulation which mimics the working of the Paxos consensus protocol. This is a re-implementation of the code in [Paxos Made Moderately Complex – van Resesse 2011](http://www.cs.cornell.edu/courses/cs7412/2011sp/paxos.pdf), which provides a deep investigation of Paxos. 

**Summary**

Paxos provides a protocol for state machine replication in a distributed asynchronous environment, that allows failures. In essence, it is used to solve a consensus problem in a distributed system. [Consensus is the process of agreeing on one result among a group of participants. This problem becomes difficult when the participants or their communication medium may experience failures](https://en.wikipedia.org/wiki/Paxos_(computer_science)).

Note: These notes primarily derived from [Paxos Made Moderately Complex](http://www.cs.cornell.edu/courses/cs7412/2011sp/paxos.pdf) 

**Processes**

The simulation contains the following participating processes:

* Client: A client process makes a request to modify or read a state. It broadcasts its request to all replica processes.
* Replica: A replica process maintains a copy of the application state. Every replica process receives requests from the clients, and asks the leaders to serialize the requests. A consistent serialization provided by this protocol allows all the replicas to see the same sequence. Every replica applies this sequence in order, to its application state.
* Leader: A leader process receives requests from the replicas. Every leader runs a two phase [SYNOD protocol](http://research.microsoft.com/en-us/um/people/lamport/pubs/lamport-paxos.pdf) with all the acceptors. The leader has two sub-processes: Scout and the Commander, which participate in phases one and two of the SYNOD protocol with the acceptor respectively.
* Acceptor: The Acceptor primarily communicates with the scout and commander and maintains its own state. Collectively, it provides the fault tolerant memory of Paxos.

To be resilient to `f` failures, we need `f+1` replicas and leaders and `2f+1` acceptors.

**Concepts**

* Command: A globally unique command initiated by a client to a replica identifying a request to modify/read the application state. Represented by a triple {client-id, command-id, operation}, where:
    * client-id is a client unique identifier.
    * command-id is a unique identifier within an individual client's space.
    * operation is the specific operation requested by the client.
* Slot: A Replica maintains an sequence of slots, which are assigned to commands (A replica "proposes" a command to a slot, a leader "decides" this proposal). This sequence of assignment should be the same across all the replicas as guaranteed by the Paxos protocol.
* Ballot, Leader maintain ballots. A ballot is a tuple of {number, leader-id}, The ballot must be lexicographically sorted, thereby requiring both the number and leader-id to be lexicographically sorted. This allows ballots to be totally ordered. Ballots are used for voting in the SYNOD protocol.
* PValue, A triple containing {ballot, slot, command} <b, s, c>, which is used to communicate ballot results from an acceptor to the leader.

**Invariants**

Invariants for the individual processes.  

Replica:
- R1: There are no two commands decided for the same slot. 
    - i.e given two Replicas r1 & r2, two commands c1 and c2. for a 
given slot s, if r1.Decisions contains (s, c1) and r2.Decisions contains (s, c2), then c1 and c2 are the same command
- R2: All commands, upto slot_out are the set of decisions.
- R3: For all replicas r, r.state is essentially the result of applying the commands in the set of decisions 
(R.Decisions) in order
- R4: For each replica, r.slot_out cannot decrease over time
- R5: A replica proposes commands for configuration it knows about. 
    - when the slot s is in *[slot_in, slot_out+WINDOW)*
    
Acceptor:
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

Commander:
- C1: For any ballot (b) slot (s) combination, atmost one Command (c) is considered. i.e atmost one Commander is spawned for
a given <b, s, c>
- C2: Suppose <b, s, c> is accepted by a majority of acceptors a. Then if a commander is spawned for a <b', s, c'> such that
b' > b, then c = c'

Notes: C1 => A4, and C2 => A5, which in turns implies R1. 
