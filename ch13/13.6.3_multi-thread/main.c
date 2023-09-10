#include <signal.h>
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>

pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;
int recieved_intr = 0;
sigset_t signals;

// シグナルハンドラではないので、非同期シグナル安全でない関数を呼び放題
void *handle_thread(void *arg)
{
  printf("handler thread waits SIGINT\n");
  int signum;
  // sigwaitを使えば、signalsに含まれるシグナルのマスクが解除される
  int err = sigwait(&signals, &signum);
  if (err || signum != SIGINT)
  {
    abort();
  }

  printf("recieved SIGINT\n");
  pthread_mutex_lock(&mutex);
  recieved_intr = 1;
  pthread_mutex_unlock(&mutex);
  return NULL;
}

int main(void)
{
  // シグナル集合を初期化
  sigemptyset(&signals);
  // シグナル集合にはSIGINTだけ入れておく
  sigaddset(&signals, SIGINT);

  pthread_t t;
  // sigint_setに含まれるシグナルをマスクする
  // この設定は、新たなスレッドにも引き継がれる
  pthread_sigmask(SIG_BLOCK, &signals, NULL);
  // スレッドを生成して実行
  pthread_create(&t, NULL, handle_thread, NULL);

  printf("waiting...\n");
  while (1)
  {
    pthread_mutex_lock(&mutex);
    if (recieved_intr)
    {
      break;
    }
    pthread_mutex_unlock(&mutex);
  }
  printf("finished\n");
  return 0;
}
