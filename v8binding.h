#ifdef __cplusplus
extern "C" {
#endif

typedef struct template_s tpl;

void init_v8();
tpl* new_template(char*);
char* eval_template(tpl*, char*);
char* get_template_err(tpl*);

#ifdef __cplusplus
} // extern "C"
#endif
