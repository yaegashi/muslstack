#include <stdio.h>
#include <pthread.h>

void *func(void *ptr)
{
        int n = *(int *)ptr + 1;
        printf("Stack level %d\n", n);
        func(&n);
        return NULL;
}

int main()
{
        int n = 0;
        pthread_t t;
        pthread_create(&t, NULL, func, &n);
        pthread_join(t, NULL);
        return 0;
}
