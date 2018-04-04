struct xx {
    int i;
    /**
	 * Text
	 */
    struct yy {
        int j;
        struct zz {
            int k;
        } deep;
    } inner;
};

int main() {
	int i = 0;
	int j = 1;
	if (i > j)
		return 1;
	if (i == j){
		return 2;
	}

	//////////////
    int value = 1;
    while (value <= 3) {
        value++;
    }

	//////////////
    switch (1) {
    case 5:
        break;
    case 2:
        break;
    }

    for (; i < 30; i++)
	{
	}

	return 0;
}

int i = 40;

void function()
{
    i += 2;
}

/* Text */
enum number { zero,
    one,
    two,
    three };

