#include <stdio.h>
#include <signal.h>
#include <stdlib.h>

#define MAX_SIZE 1024

char *info = NULL;
int signaled = 0;

void handler(int signum)
{
  printf("recieved signal %d\n", signum); // <-- アウト：I/O
  fprintf(stderr, info);                  // <-- アウト：I/O
  free(info);                             // <- アウト：mainでfree()を実行中にシグナルを受信するとメモリが壊れる可能性あり
  info = NULL;                            // <- アウト：共有オブジェクトにアクセス
  signaled = 1;                           // <- アウト：共有オブジェクトにアクセス
}

int main(void)
{
  printf("signal-unsafe start\n");

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
  free(info);
  return 0;
}
