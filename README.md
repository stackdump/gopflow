# gopflow

Go library to use PetriNets encoded using pflow schema to construct state-machines.

# Status

[![Build Status](https://travis-ci.org/stackdump/gopflow.svg?branch=master)](https://travis-ci.org/stackdump/gopflow)

Tested in Isolation - working to test within other codebases.

# Motivation

Petri-nets are well explored data structures that have mathematically verifiable properties.

States and transitions are computed as a [Vector addition System with State](https://en.wikipedia.org/wiki/Vector_addition_system)
This vector format makes machine learning analysis of event logs very trivial.

This library is compatible with `.pflow` files produced with a [visual editor](http://www.pneditor.org/)
Once a user is familiar with the basic semantics of a Petri-Net, new process flows can be developed rapidly.

