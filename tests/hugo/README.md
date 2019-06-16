# muslstack test cases using Hugo

## Successful cases

Hugo with default stack size expanded to 8MB (`--build-arg STACKSIZE=0x800000`)
always returns an expected error from hugo: `Stack depth exceeded max of 1024`.
No segmentation fault.

```console
$ docker build . --build-arg STACKSIZE=0x800000
Sending build context to Docker daemon  9.216kB
Step 1/9 : FROM alpine:edge
 ---> 43cffc6f84a4
...
...
Step 9/9 : RUN if test -n "$STACKSIZE"; then         muslstack -s "$STACKSIZE" /usr/bin/hugo;     else         muslstack /usr/bin/hugo;     fi &&     objdump -p /usr/bin/hugo | grep -A1 STACK &&     hugo || exit $?
 ---> Running in 237564d40197
/usr/bin/hugo: stackSize: 0x800000
   STACK off    0x0000000000000000 vaddr 0x0000000000000000 paddr 0x0000000000000000 align 2**4
         filesz 0x0000000000000000 memsz 0x0000000000800000 flags rw-
Building sites … ERROR 2019/06/16 14:22:30 error: failed to transform resource: SCSS processing failed: file "stdin", line 6, col 25: Stack depth exceeded max of 1024 
Error: Error building site: logged 1 error(s)
Total in 2372 ms
The command '/bin/sh -c if test -n "$STACKSIZE"; then         muslstack -s "$STACKSIZE" /usr/bin/hugo;     else         muslstack /usr/bin/hugo;     fi &&     objdump -p /usr/bin/hugo | grep -A1 STACK &&     hugo || exit $?' returned a non-zero code: 255
```

## Unsuccessful cases

You can see `Segmentation fault` in most hugo build attempts
when no stack size expansion is specified (`--build-arg STACKSIZE=`).
But sometimes it returns an expected result without segmentation fault.

```console
$ docker build . --build-arg STACKSIZE=
Sending build context to Docker daemon  10.24kB
Step 1/9 : FROM alpine:edge
 ---> 43cffc6f84a4
...
...
Step 9/9 : RUN if test -n "$STACKSIZE"; then         muslstack -s "$STACKSIZE" /usr/bin/hugo;     else         muslstack /usr/bin/hugo;     fi &&     objdump -p /usr/bin/hugo | grep -A1 STACK &&     hugo || exit $?
 ---> Running in 95dcbcd0f438
/usr/bin/hugo: stackSize: 0x0
   STACK off    0x0000000000000000 vaddr 0x0000000000000000 paddr 0x0000000000000000 align 2**4
         filesz 0x0000000000000000 memsz 0x0000000000000000 flags rw-
Building sites … Segmentation fault
The command '/bin/sh -c if test -n "$STACKSIZE"; then         muslstack -s "$STACKSIZE" /usr/bin/hugo;     else         muslstack /usr/bin/hugo;     fi &&     objdump -p /usr/bin/hugo | grep -A1 STACK &&     hugo || exit $?' returned a non-zero code: 139
```
