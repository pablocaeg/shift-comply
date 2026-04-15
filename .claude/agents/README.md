# Shift Comply Agents

Custom Claude Code agents for automating shift-comply workflows.

## new-jurisdiction

End-to-end pipeline for adding a new healthcare jurisdiction to the shift-comply database. Researches real regulations from primary legal sources, presents findings for human verification, then generates the complete Go package with tests.

### Usage

```
/agents/new-jurisdiction Add New York
/agents/new-jurisdiction Add France
/agents/new-jurisdiction Add Andalusia
```

### Pipeline

```
  Phase 1              Phase 2              Phase 3              Phase 4              Phase 5
+-----------+       +-----------+       +-----------+       +-----------+       +-----------+
| RESEARCH  |       | PRESENT   |       | GENERATE  |       | TEST      |       | VERIFY    |
|           |       |           |       |           |       |           |       |           |
| Search    |  -->  | Show      |  -->  | Write Go  |  -->  | Add unit  |  -->  | go build  |
| primary   |       | table of  |       | package   |       | tests     |       | go test   |
| legal     |       | all rules |       | matching  |       | covering  |       | go vet    |
| sources   |       | with      |       | existing  |       | registry, |       | CLI check |
|           |       | citations |       | patterns  |       | values,   |       | JSON      |
| Verify    |       | + absences|       |           |       | scopes,   |       | check     |
| every URL |       | + uncertain|      | Update    |       | chain     |       |           |
|           |       |           |       | imports   |       |           |       | Report    |
| Cross-ref |       | ASK USER  |       |           |       |           |       | results   |
| sources   |       | TO REVIEW |       |           |       |           |       |           |
+-----------+       +-----+-----+       +-----------+       +-----------+       +-----------+
                          |
                    HUMAN CHECKPOINT
                          |
                    User reviews:
                    - Every numeric value
                    - Every legal citation
                    - Every source URL
                    - Every scope assignment
                    - Every effective date
                    - Documented absences
                    - Inheritance overrides
                          |
                    Types "proceed"
```

### What it knows

The agent has deep context about:

- **The full data model**: every struct field, every JSON tag, every type
- **All available constants**: 9 categories, 40+ rule keys, 12 staff types, 17 unit types, 4 operators, 3 enforcement levels, 9 scope types, 8 units, 7 per-period types
- **How inheritance works**: child overrides parent when keys match, unmatched rules pass through
- **How scope filtering works**: empty scope = all employers, specific scope filters down
- **Common mistakes**: ACGME is not law (accreditation), Spain's Estatuto Marco is public-only, autonomous community rules are collective agreements not statutes
- **Code patterns**: date creation, boolean rules, documented absences, time-versioned values, averaging periods, exceptions
- **Style rules**: no em dashes, kebab-case keys, grouped rule functions, factual descriptions

### Quality gates

Before presenting findings, the agent verifies:
- Every value traced to a primary legal source (official gazette, legislature, legal database)
- Every URL fetched and confirmed to exist
- Scope correctly assigned (public_health for public-only, hospitals for hospital-only, etc.)
- Effective dates are actual enactment dates
- Parent jurisdiction chain is correct
- Override keys match parent exactly

### What it generates

```
jurisdictions/{code}/{code}.go    Complete Go package with all rules
comply/jurisdiction.go            Code constant added (if needed)
jurisdictions/jurisdictions.go    Blank import added
comply/comply_test.go             Tests for registry, values, scopes, chain
```

### Verified against

The generated code is built, tested, vetted, and CLI-verified before the agent reports done.
