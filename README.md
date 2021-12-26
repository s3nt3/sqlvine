# SQLVINE: Coverage-Guided Fuzzing on TiDB

Author: Zengxian Ding

# Introduction

SQLVine is a coverage-guided fuzzing framework based on golang native fuzzer for TiDB.

# Background

Ensuring the quality of TiDB's SQL layer is a complex engineering problem. The existing test cases are not enough to fully cover the SQL layer. Fuzzing is an effective way to improve test coverage. We can use fuzzing to discover new code coverage that may be missed by existing test cases. However, the existing fuzzing tools for TiDB such as [go-randgen](https://github.com/pingcap/go-randgen), [go-sqlsmith](https://github.com/PingCAP-QE/go-sqlsmith), [sql-spider](https://github.com/zyguan/sql-spider), and [go-sqlancer](https://github.com/PingCAP-QE/go-sqlancer) are all generation-based fuzzing solutions. They are very powerful in generating SQL queries, but difficult to combine with other testing techniques to extend the test coverage. To solve this problem, we designed a coverage-guided fuzzing framework based on the new feature native fuzzer which will be introduced in golang 1.18 ([Fuzzing is Beta Ready - The Go Programming Language](https://go.dev/blog/fuzz-beta)). 

# Architecture

The architecture of the system can be seen through the data flow diagram, as shown below:

![image](https://github.com/s3nt3/sqlvine/blob/main/assets/dataflow.png)
