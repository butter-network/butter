% Butter - a decentralised application framework
% Alexandre Shinebourne
% March 2022

# Outline
- An overview of the project
	- What is Butter?
	- Motivation
	- Review of the research & literature
	- Demo
- Getting technical
	- Introducing the problems
	- Butter's approach
- How we got here?
  - The path
  - Unforeseen problems
  - What's to come?

# Introducing Butter

Butter is a networking stack and framework for building decentralised applications (dapps)

![](butterLogo.png)

<!--First and formost, we have a nice logo!-->

# Motivation

During my time at university, studying Computer Science, I became wholly aware of the value of information. I think information is the most valuable thing in the world - information in its raw state. I have to come to think of information as detached from humans, even if we came up with it, it should exist as an entity in of itself. This is when I became aware of the inherent fragility of the internet - information is tightly coupled to specific hardware and hence specific people (individuals, organisation, companies...). This can lead to central points of failure, censorship.

Now lots of people are aware of this precarious state, right, e.g., effort to make decentralised currencies with the blockchain (which I discussed initially with my supervisor Adam who very quickly dismissed the Blockchain). Why? Because while it works for storing small transactions it can't be viewed as a medium for storing large parts of the internet - to compuationally demanding/message complexity demanding and too much redundancy.

This project is an effort to look at a more efficient (most notably in redundancy) decentralised approach which resulted in me building a framework for buidling applications with a unstructured peer-to-peer architecture.

# Why up2p?

Surely enough because most technologies that claim to be decentralised are not all that decentralised. During my research I came to realise that decentralisation was not a binary state but rather sliding scale. All these technologies that claim to be decenralise use Structured p2p networks with a Kadmilia Overlay network (or version of) and these require a bootsrap node - this bootstrap node could be any node but it has to be known - hence we use centralised databased of well known bootsrap nodes. This is not very fault-tolerant as i the database of bootsrap nodes goes down, then nodes that want to join the unstructured network can't find a boostrap node... This is not enough and brings us to the taxonomy of distributed systems...

# The taxonomy of distributed systems

# Demo

So now that we know what butter is and where it lies in the taxonomy of distributed systems - lets give you guys a demo

Run through the code of each example and then compile and run the demos

# Getting technical

<!--Leave blank-->

# The problems

Key problems when designing decentralised unstructured peer-to-peer architecture systems

- Discovery - cold start problem
- Known host management - heavily AI inspired (GOFAI)
- overlay network for persistent information 
  - information retrieval - heavily AI inspired (GOFAI)
  - persistent fault-tolerant storage (high availability)
- NAT traversal

# Butter's approach

 how the framework addresses them

| Problem name          | Butter's solution                              |
| --------------------- | ---------------------------------------------- |
| Discovery             | Multicast                                      |
| Known host management | AI inspired (GOFAI) - state and value function |
| IR                    | AI inspired (GOFAI) - directed BFS             |
| Persistent storage    | PCG overlay                                    |
| NAT Traversal         | (Imperfect) Ambassador nodes                   |

Butter is greedy - it will try and use as much memory as you allow it to use

# Discovery protocol

# Known host management protocol

# How we got here

<!--Add project management story I.e. found progress to be very slow in rust and wasn’t the right tool for the job so made the decision to switch to go and project improved-->
