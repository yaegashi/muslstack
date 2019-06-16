# Confirming stack size expansion by muslstack

You can confirm stack levels it reaches
before hitting segmentation fault at the bottom of the stack.

```console
$ docker build . --build-arg STACKSIZE=
...
Step 9/9 : RUN if test -n "$STACKSIZE"; then         muslstack -s "$STACKSIZE" a.out;     else         muslstack a.out;     fi &&     objdump -p a.out | grep -A1 STACK &&     ./a.out >log ||     echo -e "\nExit $?" >>log && tail -5 log && exit 
1
 ---> Running in 9f47b6a83d86
a.out: stackSize: 0x0
   STACK off    0x0000000000000000 vaddr 0x0000000000000000 paddr 0x0000000000000000 align 2**4
         filesz 0x0000000000000000 memsz 0x0000000000000000 flags rw-
Segmentation fault
Stack level 2726
Stack level 2727
Stack level 2728

Exit 139
The command '/bin/sh -c if test -n "$STACKSIZE"; then         muslstack -s "$STACKSIZE" a.out;     else         muslstack a.out;     fi &&     objdump -p a.out | grep -A1 STACK &&     ./a.out >log ||     echo -e "\nExit $?" >>log && tail -5 log && exit 1' returned a non-zero code: 1
```

```console
$ docker build . --build-arg STACKSIZE=0x800000
...
Step 9/9 : RUN if test -n "$STACKSIZE"; then         muslstack -s "$STACKSIZE" a.out;     else         muslstack a.out;     fi &&     objdump -p a.out | grep -A1 STACK &&     ./a.out >log ||     echo -e "\nExit $?" >>log && tail -5 log && exit 1
 ---> Running in b70837f6fab5
a.out: stackSize: 0x800000
   STACK off    0x0000000000000000 vaddr 0x0000000000000000 paddr 0x0000000000000000 align 2**4
         filesz 0x0000000000000000 memsz 0x0000000000800000 flags rw-
Segmentation fault
Stack level 174779
Stack level 174780
Stack level 174781
Stack level 174782
Exit 139
The command '/bin/sh -c if test -n "$STACKSIZE"; then         muslstack -s "$STACKSIZE" a.out;     else         muslstack a.out;     fi &&     objdump -p a.out | grep -A1 STACK &&     ./a.out >log ||     echo -e "\nExit $?" >>log && tail -5 log && exit 1' returned a non-zero code: 1
```