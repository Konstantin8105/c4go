static char *sstr_s;
char sstr_bufs[10];
int sstr_n;

char *sstr_pop(void)
{
	char *ret = sstr_s;
	sstr_s = sstr_bufs[--sstr_n];
	return ret;
}

int sstr_next(void)
{
	return *sstr_s ? (unsigned char) *sstr_s++ : -1;
}

void sstr_back(int c)
{
	sstr_s--;
}
