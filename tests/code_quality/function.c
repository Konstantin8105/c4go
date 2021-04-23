static char *sstr_s;
char sstr_bufs[10];
int sstr_n;

char *sstr_pop(void)
{
	char *ret = sstr_s;
	sstr_s = &sstr_bufs[--sstr_n];
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

int sf1() {
	sstr_n++;
	return *sstr_s == sstr_bufs[sstr_n] ? 1: 0;
}

int st2() {
	sstr_n++;
	int s = *sstr_s == sstr_bufs[sstr_n] ? 1: 0;
	sstr_n--;
	return s;
}

int st3() {
	sstr_n++;
	int s = (*sstr_s == sstr_bufs[sstr_n] ? sstr_n+1: sstr_n-1);
	sstr_n--;
	return s;
}

int st4() {
	return (*sstr_s == sstr_bufs[sstr_n] ? sstr_n+1: sstr_n-1);
}
