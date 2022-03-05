

# Motivation

During my time at university, studying Computer Science, I became wholly aware of the value of information. I think information is the most valuable thing in the world - information in its raw state. I have to come to think of information as detached from humans, even if we came up with it, it should exist as an entity in of itself. This is when I became aware of the inherent fragility of the internet - information is tightly coupled to specific hardware and hence specific people (individuals, organisation, companies...). This can lead to central points of failure, censorship.

Now lots of people are aware of this precarious state, right, e.g., effort to make decentralised currencies with the blockchain (which I discussed initially with my supervisor Adam who very quickly dismissed the Blockchain). Why? Because while it works for storing small transactions it can't be viewed as a medium for storing large parts of the internet - to compuationally demanding/message complexity demanding and too much redundancy.

This project is an effort to look at a more efficient (most notably in redundancy) decentralised approach which resulted in me building a framework for buidling applications with a unstructured peer-to-peer architecture.

# Why up2p?

Surely enough because most technologies that claim to be decentralised are not all that decentralised. During my research I came to realise that decentralisation was not a binary state but rather sliding scale. All these technologies that claim to be decenralise use Structured p2p networks with a Kadmilia Overlay network (or version of) and these require a bootsrap node - this bootstrap node could be any node but it has to be known - hence we use centralised databased of well known bootsrap nodes. This is not very fault-tolerant as i the database of bootsrap nodes goes down, then nodes that want to join the unstructured network can't find a boostrap node... This is not enough and brings us to the taxonomy of distributed systems...


