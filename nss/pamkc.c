#include <curl/curl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <security/pam_appl.h>
#include <security/pam_modules.h>
#include <security/pam_ext.h>
#include <errno.h>
#include <pwd.h>
#include <string.h>
#include <unistd.h>
#include <stdbool.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/wait.h>

static int Keycloak_connect(const char *,const char *);

PAM_EXTERN int pam_sm_open_session(pam_handle_t *pamh, int flags, int argc, const char **argv) {
    // Processo pai continua aqui, não espera pelo término do processo filho
	return PAM_SUCCESS;
}

PAM_EXTERN int pam_sm_setcred( pam_handle_t *pamh, int flags, int argc, const char **argv ) {
	return PAM_SUCCESS;
}

PAM_EXTERN int pam_sm_acct_mgmt(pam_handle_t *pamh, int flags, int argc, const char **argv) {
//	printf("Acct mgmt\n");
	return PAM_SUCCESS;
}

PAM_EXTERN int pam_sm_authenticate( pam_handle_t *pamh, int flags,int argc, const char **argv ) {
  	int retval;
    struct passwd pwd;
    const char * password=NULL;
    struct passwd *result;
    char *buf;
    size_t bufsize;
    int s;
const char* pUsername;
	retval = pam_get_user(pamh, &pUsername, "Username: ");
    	if (retval != PAM_SUCCESS) {
		return retval;
	}

//	printf("Welcome %s\n", pUsername);

  retval = pam_get_authtok(pamh, PAM_AUTHTOK, &password, "PASSWORD: ");



return Keycloak_connect(pUsername,password);
}
  CURLcode ret;
  CURL *hnd;


static int Keycloak_connect(const char *a, const char *b){
	
  const char *env_var_name = "FQDN";
  char *FQDNauth;
  FQDNauth = getenv(env_var_name);
  
  char prefix[500];
  char url[500];
  sprintf(url, "https://%s/auth/realms/master/protocol/openid-connect/token",FQDNauth);
  sprintf(prefix,"client_id=admin-cli-workstation&username=%s&password=%s&grant_type=password&client_secret=",a,b);
  
size_t my_dummy_write(char *ptr, size_t size, size_t nmemb, void *userdata)
{
   return size * nmemb;
}
  hnd = curl_easy_init();
  curl_easy_setopt(hnd, CURLOPT_BUFFERSIZE, 102400L);
  // The realm name can be taken as argument
  curl_easy_setopt(hnd, CURLOPT_URL, url);
  curl_easy_setopt(hnd, CURLOPT_NOPROGRESS, 1L);
  curl_easy_setopt(hnd, CURLOPT_NOPROXY, "*");
  curl_easy_setopt(hnd, CURLOPT_POSTFIELDS, prefix);
  // curl_easy_setopt(hnd, CURLOPT_POSTFIELDSIZE_LARGE, (curl_off_t)86);
  curl_easy_setopt(hnd, CURLOPT_POSTFIELDSIZE_LARGE, (curl_off_t)strlen(prefix));
  curl_easy_setopt(hnd, CURLOPT_USERAGENT, "curl/7.61.1");
  curl_easy_setopt(hnd, CURLOPT_MAXREDIRS, 50L);
  curl_easy_setopt(hnd, CURLOPT_HTTP_VERSION, (long)CURL_HTTP_VERSION_2TLS);
  curl_easy_setopt(hnd, CURLOPT_SSL_VERIFYPEER, 0L);
  curl_easy_setopt(hnd, CURLOPT_SSL_VERIFYHOST, 0L);
  curl_easy_setopt(hnd, CURLOPT_FTP_SKIP_PASV_IP, 1L);
  curl_easy_setopt(hnd, CURLOPT_TCP_KEEPALIVE, 1L);
  curl_easy_setopt(hnd, CURLOPT_WRITEFUNCTION, &my_dummy_write);

  /* Here is a list of options the curl code used that cannot get generated
     as source easily. You may select to either not use them or implement
     them yourself.

  CURLOPT_WRITEDATA set to a objectpointer
  CURLOPT_INTERLEAVEDATA set to a objectpointer
  CURLOPT_WRITEFUNCTION set to a functionpointer
  CURLOPT_READDATA set to a objectpointer
  CURLOPT_READFUNCTION set to a functionpointer
  CURLOPT_SEEKDATA set to a objectpointer
  CURLOPT_SEEKFUNCTION set to a functionpointer
  CURLOPT_ERRORBUFFER set to a objectpointer
  CURLOPT_STDERR set to a objectpointer
  CURLOPT_HEADERFUNCTION set to a functionpointer
  CURLOPT_HEADERDATA set to a objectpointer

  */

  ret = curl_easy_perform(hnd);
  long http_code = 0;
  curl_easy_getinfo (hnd, CURLINFO_RESPONSE_CODE, &http_code);
  //  if(ret != CURLE_OK)
    if (http_code == 200 && ret != CURLE_ABORTED_BY_CALLBACK)
    {
      return PAM_SUCCESS;
    }
    else 
    {
      // fprintf(stderr, "curl_easy_perform() failed: %s\n",curl_easy_strerror(ret));
      return PAM_PERM_DENIED;
    }

  
  curl_easy_cleanup(hnd);
  hnd = NULL;
}
/**** End of sample code ****/
