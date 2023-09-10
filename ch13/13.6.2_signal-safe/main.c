#include <stdio.h>
#include <signal.h>
#include <stdlib.h>

#define MAX_SIZE 1024

char *info = NULL;
volatile sig_atomic_t signaled = 0;

void handler(int signum)
{
  signaled = 1; // <-- OK：sig_atomic_t型のグローバル変数へのアクセス
}

int main(void)
{
  printf("signal-safe start\n");

  // シグナルハンドラ登録
  if (signal(SIGINT, handler) == SIG_ERR)
  {
    printf("SIG_ERR: SIGINT\n");
    exit(1);
  }
  info = (char *)malloc(MAX_SIZE);
  if (info == NULL)
  {
    printf("malloc error\n");
    exit(1);
  }

  printf("waiting...\n");
  while (!signaled)
  {
    fprintf(stderr, info);
  }

  printf("recieved SIGINT\n");
  fprintf(stderr, info);
  free(info);
  info = NULL;

  return 0;
}
