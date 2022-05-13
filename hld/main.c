
#include <event2/listener.h>
#include <event2/thread.h>
#include <event2/bufferevent.h>
#include <event2/buffer.h>
#include <sys/sendfile.h>

#include <arpa/inet.h>

#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <pthread.h>
#include <signal.h>
#include <stdbool.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <unistd.h>

#define IS_EXT(x, ext) (strncasecmp(x, ext, sizeof(ext)) == 0)
#define JOB_QUEUE_LEN 4096
#define BUFF_LEN 2048

const int port = 8080;
const int RETRIES = 3;
const int THREAD_NUM = 4;
const useconds_t TIMEOUT_INTERVAL = 10;
const char RESPONSE_OK_FMT[] = "HTTP/1.1 200 OK\r\nConnection: close\r\nServer: Ktlh/1.0 (Uwuntu)\r\nContent-type: %s\r\nContent-Length: %zd\r\n\r\n";
const char FORBIDDEN[] = "HTTP/1.1 403 Forbidden\r\nServer: Ktlh/1.0 (Uwuntu)\r\nConnection: close\r\n\r\n";
const char NOT_FOUND[] = "HTTP/1.1 404 Not Found\r\nServer: Ktlh/1.0 (Uwuntu)\r\nConnection: close\r\n\r\n<html><body><h1>404 Not Found</h1></body></html>";
const char METHOD_UNIMPLEMENTED[] = "HTTP/1.1 405 Method Not Implemented\r\nServer: Ktlh/1.0 (Uwuntu)\r\nConnection: close\r\n\r\n";
const char TOO_LONG[] = "HTTP/1.1 414 Request-URI Too Long\r\nServer: Ktlh/1.0 (Uwuntu)\r\nConnection: close\r\n\r\n";
const char INDEX[] = "/index.html";
const char wd[] = "/home/kali/hld/www/";

pthread_mutex_t mutex;
int queue[JOB_QUEUE_LEN];
int* head;
int* tail;

void enqueue (int n){
	while (1){
		if(*tail == -1){
			pthread_mutex_lock(&mutex);
			if (*tail == -1){
				*tail = n;
				tail = tail-queue >= JOB_QUEUE_LEN - 1 ? queue : tail+1;
				pthread_mutex_unlock(&mutex);
				return;
			}
			pthread_mutex_unlock(&mutex);
		}else{
			usleep(0);
		}
	}
}



int dequeue (){
	int n;
	while (1){
		if(*head != -1){
			pthread_mutex_lock(&mutex);
			if (*head != -1){
				n = *head;
				*head = -1;
				head = head-queue >= JOB_QUEUE_LEN - 1 ? queue : head+1;
				pthread_mutex_unlock(&mutex);
				return n;
			}
			pthread_mutex_unlock(&mutex);
		}else{
			usleep(0);
		}
	}
}
		



char hex2char (char c)
{
    if ('0' <= c && c <= '9') return c - '0';
    if ('A' <= c && c <= 'F') return c - 'A' + 10;
    if ('a' <= c && c <= 'f') return c - 'a' + 10;
    return -1;
}
static int inpl_urldecode(char* string){
	int len = strlen(string);
	char* cur_pos = string;
	for (int i = 0; i<len; i++){
		if (string[i] == '%'){
			if ((i+2 >= len)||(hex2char(string[i+1]) == -1)||(hex2char(string[i+2]) == -1)) {
				return -1;
			}
			*(cur_pos++) = (hex2char(string[i+1])<<4)+hex2char(string[i+2]);
			i+=2;
		}else if (string[i] == '+'){
			*(cur_pos++) = ' ';
		}else{
			*(cur_pos++) = string[i];
		}
	}
	*cur_pos=0;
	return 0;
}

static const char* resolve_content_type(char* file_extention){
	if (IS_EXT(file_extention, ".html")) {
		return "text/html";
	} else if (IS_EXT(file_extention, ".js")) {
		return "application/javascript";
	} else if (IS_EXT(file_extention, ".css")) {
		return "text/css";
	} else if (IS_EXT(file_extention, ".jpg")||IS_EXT(file_extention, ".jpeg")) {
		return "image/jpeg";
	} else if (IS_EXT(file_extention, ".png")) {
		return "image/png";
	} else if (IS_EXT(file_extention, ".gif")) {
		return "image/gif";
	} else if (IS_EXT(file_extention, ".swf")) {
		return "application/x-shockwave-flash";
	}
	return "binary/octet-stream";
}


static void accept_conn_cb(struct evconnlistener *listener, evutil_socket_t sock,
		struct sockaddr *address, int socklen,
		void *ctx){
	enqueue(sock);
}




void* worker_func(void* opaque){
	char buff[BUFF_LEN];
	char response_ok[BUFF_LEN];
	char *url;
	char *method;
	int res;
	bool send_file;
	evutil_socket_t sock;
	struct job *j;
	while (1){
		sock = dequeue();
		res = -1;
		send_file = true;
	//	for (int i = 0; i < RETRIES && res == -1 && errno == EAGAIN; res = recv(sock, buff, sizeof(buff), 0), i++){
	//		usleep(TIMEOUT_INTERVAL * i);
	//	}
		res = recv(sock, buff, sizeof(buff), 0);
		if ( res > 0 ){
			method = strtok(buff, "? \r");
			if (method && (!strncmp(method,"GET",3) || !strncmp(method,"HEAD",4) ) ){
				url = strtok(NULL, "? \r");
				if (url && url[0]=='/' && strtok(NULL,"\r")){
					struct stat file_stat;
					url = url - 1;
					url[0] = '.';
					inpl_urldecode(url);
					int url_len = strlen(url);
					int stat_res;
					if (url[url_len-1] == '/'){
						url[url_len-1] = 0;
						stat_res = stat(url, &file_stat);
						if (S_ISDIR(file_stat.st_mode)){
							strncat(url, INDEX, sizeof(INDEX));
							stat_res = stat(url, &file_stat);
						} else{
							stat_res = -1;
						}
					} else {
						stat_res = stat(url, &file_stat);
					}
					if (strstr(url, "/../")){
						send(sock, FORBIDDEN, sizeof(FORBIDDEN), 0);
					}else if (!stat_res){
						snprintf(response_ok, sizeof(response_ok)-1, RESPONSE_OK_FMT, resolve_content_type(strrchr(url,'.')), file_stat.st_size);
						send(sock, response_ok, strlen(response_ok), 0);
						if (!strncmp(method,"GET",3)){
							int fd;
							fd = open(url, O_RDONLY);
							sendfile(sock, fd, 0, file_stat.st_size);
							close(fd);
						}
					}else if ((url[url_len-1] == '/') ){
						send(sock, FORBIDDEN, sizeof(FORBIDDEN), 0);
					}else{
						int e = errno;
						send(sock, NOT_FOUND, sizeof(NOT_FOUND), 0);
					}
				}
			} else {
				send(sock, METHOD_UNIMPLEMENTED, sizeof(METHOD_UNIMPLEMENTED), 0);
			}
		}
		close(sock);
	}
}
void exit_on_term(int i){
	dprintf(2,"exiting\n");
	exit(0);
}

int main(int argc, char** argv){
	int sc_listen;
	struct sockaddr_in sin;
	struct event_base* evbase;
        struct evconnlistener *listener;
	pthread_t threads[THREAD_NUM];

	int reuse = 1;
	chdir(wd);
	evthread_use_pthreads();
	evbase = event_base_new();

	signal(SIGPIPE, SIG_IGN);
	signal(SIGTERM, &exit_on_term);
	memset(&sin, 0, sizeof(sin));
	sin.sin_family = AF_INET;
	sin.sin_addr.s_addr = htonl(0);
	sin.sin_port = htons(port);
	
	for (int i = 0; i < JOB_QUEUE_LEN; i++){
		queue[i] = -1;
	}
	head = queue;
	tail = queue;
	pthread_mutex_init(&mutex, NULL);
	for (int i = 0; i < THREAD_NUM; i++){
		pthread_create(&threads[i], NULL, &worker_func, NULL);
	}


	listener = evconnlistener_new_bind(evbase, accept_conn_cb, NULL,
//            LEV_OPT_CLOSE_ON_FREE|LEV_OPT_REUSEABLE|LEV_OPT_LEAVE_SOCKETS_BLOCKING, -1,
            LEV_OPT_REUSEABLE|LEV_OPT_LEAVE_SOCKETS_BLOCKING, -1,
            (struct sockaddr*)&sin, sizeof(sin));
	if (!listener) {
                perror("Couldn't create listener");
                exit(1);
        }


	event_base_dispatch(evbase);
	
	//!TODO free stale events, buffers, sockets, etc.
	event_base_free(evbase);
	return 0;
}



