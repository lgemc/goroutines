# Patrones de concurrencia en otros lenguajes

El siguiente ejemplo muestra como se implementa en diferentes lenguajes el siguiente ejemplo
```
Crea dos procesos concurrentes que juegan pin pong entre ellos

Prompt usado pa GPT 🌞


Give me an example of create two threads or process, and communicate them using classic ping-pong example

in:

golang
c
python
java
erlang
scala
```

## Golang
```go
package main

import (
    "fmt"
    "time"
)

func player1(ball chan int, done chan bool) {
    for {
        msg := <-ball
        fmt.Println("Player 1 received:", msg)
        time.Sleep(time.Second)
        msg++
        ball <- msg
        if msg >= 10 {
            done <- true
            return
        }
    }
}

func player2(ball chan int, done chan bool) {
    for {
        msg := <-ball
        fmt.Println("Player 2 received:", msg)
        time.Sleep(time.Second)
        msg++
        ball <- msg
        if msg >= 10 {
            done <- true
            return
        }
    }
}

func main() {
    ball := make(chan int)
    done := make(chan bool)

    go player1(ball, done)
    go player2(ball, done)

    ball <- 0

    <-done
}
```


# C (con pthreads)
```c
#include <stdio.h>
#include <pthread.h>
#include <unistd.h>

#define MAX_COUNT 10

pthread_mutex_t lock;
int ball = 0;

void* player1(void* arg) {
    while (1) {
        pthread_mutex_lock(&lock);
        printf("Player 1 received: %d\n", ball++);
        pthread_mutex_unlock(&lock);
        sleep(1);
        if (ball >= MAX_COUNT) {
            pthread_exit(NULL);
        }
    }
}

void* player2(void* arg) {
    while (1) {
        pthread_mutex_lock(&lock);
        printf("Player 2 received: %d\n", ball++);
        pthread_mutex_unlock(&lock);
        sleep(1);
        if (ball >= MAX_COUNT) {
            pthread_exit(NULL);
        }
    }
}

int main() {
    pthread_t t1, t2;
    pthread_mutex_init(&lock, NULL);

    pthread_create(&t1, NULL, player1, NULL);
    pthread_create(&t2, NULL, player2, NULL);

    pthread_join(t1, NULL);
    pthread_join(t2, NULL);

    pthread_mutex_destroy(&lock);
    return 0;
}
```

# Python
```python
import threading
import time

def player1(ball, lock):
    while True:
        with lock:
            print("Player 1 received:", ball)
            ball += 1
            time.sleep(1)
            if ball >= 10:
                return

def player2(ball, lock):
    while True:
        with lock:
            print("Player 2 received:", ball)
            ball += 1
            time.sleep(1)
            if ball >= 10:
                return

def main():
    ball = 0
    lock = threading.Lock()

    t1 = threading.Thread(target=player1, args=(ball, lock))
    t2 = threading.Thread(target=player2, args=(ball, lock))

    t1.start()
    t2.start()

    t1.join()
    t2.join()

if __name__ == "__main__":
    main()
```

# Java
```java
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

public class Main {
    public static void main(String[] args) {
        Lock lock = new ReentrantLock();
        int ball = 0;

        Thread t1 = new Thread(() -> {
            while (true) {
                lock.lock();
                System.out.println("Player 1 received: " + ball);
                ball++;
                lock.unlock();
                try {
                    Thread.sleep(1000);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
                if (ball >= 10) {
                    return;
                }
            }
        });

        Thread t2 = new Thread(() -> {
            while (true) {
                lock.lock();
                System.out.println("Player 2 received: " + ball);
                ball++;
                lock.unlock();
                try {
                    Thread.sleep(1000);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
                if (ball >= 10) {
                    return;
                }
            }
        });

        t1.start();
        t2.start();

        try {
            t1.join();
            t2.join();
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }
}
```

# Erlang
```erlang
-module(ping_pong).
-export([start/0, player1/1, player2/1]).

player1(Ball) ->
    receive
        {ping, Sender} ->
            io:format("Player 1 received: ~w~n", [Ball]),
            NewBall = Ball + 1,
            Sender ! {pong, NewBall},
            player1(NewBall)
    end.

player2(Ball) ->
    receive
        {pong, Sender} ->
            io:format("Player 2 received: ~w~n", [Ball]),
            NewBall = Ball + 1,
            Sender ! {ping, self()},
            player2(NewBall)
    end.

start() ->
    spawn(ping_pong, player1, [0]),
    spawn(ping_pong, player2, [0]),
    ping_pong:player1(0).
```

# Scala
```scala
import akka.actor._

case class Ping(ball: Int)
case class Pong(ball: Int)

class Player1 extends Actor {
  def receive = {
    case Ping(ball) =>
      println(s"Player 1 received: $ball")
      val newBall = ball + 1
      sender() ! Pong(newBall)
  }
}

class Player2 extends Actor {
  def receive = {
    case Pong(ball) =>
      println(s"Player 2 received: $ball")
      val newBall = ball + 1
      sender() ! Ping(newBall)
  }
}

object PingPong extends App {
  val system = ActorSystem("PingPongSystem")
  val player1 = system.actorOf(Props[Player1], name = "player1")
  val player2 = system.actorOf(Props[Player2], name = "player2")

  player1 ! Ping(0)
}
```