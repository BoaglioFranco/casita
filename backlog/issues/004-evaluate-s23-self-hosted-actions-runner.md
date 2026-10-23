---
id: ISSUE-004
status: open
priority: low
type: research
created: 2026-10-23
files:
  - .github/workflows/
---

# Evaluate a Galaxy S23 Ultra as a self-hosted GitHub Actions runner

## Problem

Investigate whether a Samsung Galaxy S23 Ultra can host a full ARM64 Linux
virtual machine used as an experimental self-hosted GitHub Actions runner. The
GitHub runner and build tools should run inside the VM rather than directly in
Android or Termux. The goal is to learn how self-hosted runners work and
explore reuse of ARM64 phone hardware, not to replace the dependable production
build path.

## Evidence

- The S23 Ultra uses an ARM64 processor, and GitHub publishes an ARM64 Linux
  runner application.
- GitHub officially supports the runner on specific Linux distributions,
  Windows, and macOS; Android and Termux are not listed as supported operating
  systems. A supported ARM64 Linux distribution inside a VM would more closely
  match GitHub's documented environment:
  https://docs.github.com/en/actions/reference/runners/self-hosted-runners
- GitHub requires Linux and Docker for workflows that use Docker container
  actions or service containers. The VM approach should be evaluated for a
  functional Linux kernel, hardware virtualization, networking, and Docker
  support; the phone and Android version may restrict some of these features.
- Android background-process restrictions, thermal throttling, battery health,
  storage, and network availability may affect runner reliability.
- A self-hosted runner executes workflow code on the device, so repository
  trust and isolation are material security concerns on a personal phone.

## Expected outcome

Determine whether the phone can safely host a usable ARM64 Linux VM, execute a
small repository-scoped workflow inside it, document the setup and removal
process, and decide whether it is useful for native ARM64 experiments. Keep
normal production builds independent of the phone.

## Acceptance criteria

- [ ] Available hardware-accelerated and emulated VM approaches are compared
- [ ] Root requirements and Android-version constraints are documented
- [ ] A supported ARM64 Linux distribution boots with persistent storage and networking
- [ ] The GitHub runner starts inside the Linux VM and appears online
- [ ] A repository-scoped workflow successfully runs a harmless diagnostic job
- [ ] Docker or an alternative container-image build path is evaluated explicitly
- [ ] Secrets, workflow trust, device permissions, and runner isolation are documented
- [ ] Battery, thermals, background execution, storage, and network behavior are tested
- [ ] Runner labels prevent normal production jobs from accidentally targeting the phone
- [ ] Runner removal and credential cleanup are documented and verified

## Notes

Treat this as a learning experiment. Do not expose the runner to workflows from
untrusted pull requests, and do not store production deployment secrets on the
phone until the security model has been reviewed.
