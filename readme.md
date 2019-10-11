## Freenet

Team: Victor Ongkowijaya (vo6@princeton.edu), Satadal Sengupta
(satadal.sengupta@cs.princeton.edu)
 
## Proposal

### Background
 
- _Name of paper_: [Freenet: A Distributed Anonymous Information Storage and
Retrieval System, by Clarke et al.](http://homepage.divms.uiowa.edu/~ghosh/freenet.pdf)

- _Brief summary of paper’s problem domain / challenge, goals, and technical
approach_: Distributed file systems provide little privacy for both producers
and consumers of data. Freenet’s goal is to solve this issue by achieving a
number of properties: anonymity for all users, deniability for storers of data,
decentralization, and resistance to censorship. Clarke et al. approached this
issue by designing a system which utilizes encryption and signing schemes
between cooperative nodes, along with an adaptive routing algorithm. 

- _Summary of paper’s current implementation, evaluation strategy, and key
results_: The system comprises of user nodes which shares their unused disk
space. This space is used to store files as well as a dynamic routing table.
Freenet utilizes three types of keys to achieve its desired properties: the
keyword signed key (KSK), signed subspace key (SSK), and content hash key
(CHK). These keys work together for nodes to identify, encrypt, retrieve, and
insert files. Nodes can find each other by a routing protocol similar to IP,
where key requests are forwarded from node to node in a chain of proxy
requests. Each node then consults its dynamic routing table to determine the
forwarding target. Clarke et al. evaluated Freenet’s efficiency (request
pathlength), scalability (request pathlength vs. number of nodes), and fault
tolerance (request pathlength as nodes failed). They also evaluate their
system’s security, providing a discussion of possible attacks and security
measures.

### Plan
 
- _Proposed implementation (language, framework, etc.)_: We plan to use Go to
re-implement the code to run on each node, which includes various aspects: key
scheme, bootstrapping, node addition/removal, fault tolerance, and routing
protocol. To simplify our evaluations, a user-friendly interface to simulate a
user running a node might be useful.

- _Evaluation strategy (testing platform/setup, simulated data/traces, etc.)_:
The paper simulated Freenet using a test network of 1000 nodes, each with a
datastore capacity of 50 items and a routing table size of 250 addresses, in a
standard ring-lattice topology. The simulated traffic consists of random
inserts interspersed with random requests with a hops-to-live value of 20. We
plan to replicate this, with potential additions outlined in section 5 below.

- _Key results trying to reproduce_: The most important results to reproduce will
be the efficiency evaluation and scalability evaluation. We want to see that
the request pathlength decreases rapidly over time as the system is used,
converging to a small median (6 in the paper). We also want to see that the
system scale to a very large number of nodes while maintaining an acceptable
efficiency (1 million nodes, 30 median pathlength).

- _Discussion of how you can compare your findings (quantitatively,
qualitatively) with previously published results_: We can construct the same
types of graphs and compare the results. The paper does not provide much
quantitative evaluation, which could be a potential addition for our project.
We can evaluate from the perspective of a user (as opposed to the existing
study which utilizes random test traffic) and then document the efficiency and
user experience, which should yield more realistic results.

- _New questions/settings trying to evaluate that are not addressed in the
original paper_: It would be interesting to evaluate the system from an
adversarial perspective and try to compromise the anonymity and deniability
guarantees. There are also a few evaluations which would provide new insight,
such as in fault tolerance (how many requests were failed), efficiency (actual
experienced latency as opposed to number of hops), persistence (how long files
persist and its relation to usage), and topology (behavior of system in more
real-world accurate network topologies) among many potential others.
