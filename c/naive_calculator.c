// sean at shanghai
// 2070311
// caculator of add/minus/multiply/divide

#include <stdio.h>
#include <stdlib.h>


// operators
const int SIGN_BAD = 0;
const int SIGN_ADD = 1;
const int SIGN_MINUS = 2;
const int SIGN_MUL = 3;
const int SIGN_DIV = 4;


// errors
const int E_BAD_SIGN = 1; // bad operator, like ^
const int E_BAD_NUM = 2; // bad num
const int E_BAD_MEM = 3; // no memory
const int E_BAD_MUL_DIV = 4; // bad multiply and divide


// linked list of calculator
struct calc {
	// current value
	float val;
	// next operation
	struct calc *next;
	// next level, like +- to */
	struct calc *next_level;
	// operator, like +-*/
	int sign;
	// if we have next_level, here is the sign after next_level expression
	int next_sign;
};


// helper functions
// set all calc members to 0
void
init_calc(struct calc *c) {
	if (c == NULL) {
		return;
	}
	c->val = 0.0;
	c->next = NULL;
	c->next_level = NULL;
	c->sign = 0;
	c->next_sign = 0;
}

void
print_calc(struct calc *c) {
	if (c == NULL) {
		printf("NULL\n");
		return;
	}
	printf("---------------------------\n");
	printf("calc.val %f\n", c->val);
	printf("calc.sign %d\n", c->sign);
	printf("calc.next_sign %d\n", c->next_sign);
	printf("calc.next_level:");
	print_calc(c->next_level);
	printf("calc.next:");
	print_calc(c->next);
	return;
}


// read an operator
// only support + - * and /
// return the offset of the expression
int
read_sign(int *sign, const char *exp) {
	printf("read _sign %s\n", exp);
	switch (exp[0]) {
		case '+':
			*sign = SIGN_ADD;
			return 1;
		case '-':
			*sign = SIGN_MINUS;
			return 1;
		case '*':
			*sign = SIGN_MUL;
			return 1;
		case '/':
			*sign = SIGN_DIV;
			return 1;
		default:
			printf("bad sign %d\n", exp[0]);
			return -E_BAD_SIGN;
	}
	return -E_BAD_SIGN;
}


// atoi atof
int
read_num(float *num, const char *exp) {
	float tmp = 0;
	int dig, dot = 0, div, offset = 0;
	printf("read_num %s\n", exp);
	if (exp[offset] > '9' || exp[offset] < '0') {
		return -E_BAD_NUM;
	}
	while(1) {
		if (exp[offset] == '.') {
			dot = 1;
			div = 10;
			offset ++;
			continue;
		}
		if (exp[offset] > '9' || exp[offset] < '0') {
			*num = tmp;
			break;
		}
		dig = exp[offset] - '0';
		printf("dig is %d\n", dig);
		if (dot == 1) {
			tmp = tmp + dig / div;
			div = div * 10;
		} else {
			tmp = tmp * 10 + dig;
		}
		offset ++;
	}
	return offset;
}


int
parse_mul_div(struct calc *c, const char *exp) {
	int sign, ret, offset = 0;
	float num;
	printf("parse mul div %s\n", exp);
	while (1) {
		if (exp[0] == 0) {
			break;
		}
		ret = read_num(&num, exp + offset);
		if (ret < 0) {
			return ret;
		}
		offset += ret;
		if (exp[offset] == 0) {
			break;
		}
		ret = read_sign(&sign, exp + offset);
		if (ret < 0) {
			return ret;
		}
		offset += ret;
		if (sign != SIGN_MUL && sign != SIGN_DIV) {
			return -E_BAD_SIGN;
		}
	}
	return offset;
}

int
read_mul_div(struct calc *p, const char *exp) {
	int sign, ret, offset;
	struct calc tmp;
	ret = parse_mul_div(p, exp);
	if (ret < 0) {
		return ret;
	}
	offset = ret;
	p->next_level = tmp.next;
	p->next = NULL;
	p->next_sign = tmp.sign;
	ret = read_sign(&sign, exp);
	if (ret < 0) {
		return ret;
	}
	offset = offset + ret;
	p->sign = sign;
	return offset;
}


struct calc*
new_calc(struct calc *prev, float val, int sign) {
	struct calc *append;
	append = (struct calc*)malloc(sizeof(struct calc));
	if (append == NULL) {
		return append;
	}
	init_calc(append);
	append->val = val;
	append->sign = sign;
	if (prev == NULL) {
		return append;
	}
	prev->next = append;
	return append;
}


int
parse_add_minus(struct calc **c, const char *exp) {
	int num_len, sign, ret, offset = 0;
	float num;
	struct calc *p = NULL;
	for (;;) {
		ret = read_num(&num, exp + offset);
		if (ret < 0) {
			return 0;
		}
		printf("read num %f\n", num);
		offset = offset + ret;
		num_len = ret;
		if (exp[offset] == 0) {
			break;
		}
		ret = read_sign(&sign, exp + offset);
		if (ret < 0) {
			return 0;
		}
		printf("read sign %d\n", sign);
		offset += ret;
		if (sign == SIGN_MUL || sign == SIGN_DIV) {
			// drop back number and sign
			offset = offset - 1 - num_len;
			ret = parse_mul_div(p, exp + offset);
			if (ret < 0) {
				return -E_BAD_MUL_DIV;
			}
			offset += ret;
		}
		p = new_calc(p, num, sign);
	}
	p = new_calc(p, num, SIGN_BAD);
	*c = p;
	return 0;
}


int
calc_add_minus(struct calc *head, float *ret) {
	float tmp = 0;
	struct calc *p = head;
	for {
		if (p->next == NULL && p->next_level == NULL) {
			break;
		}
	}
	*ret = tmp;
	return 0;
}


int
main(int argc, char **argv) {
	if (argc < 2) {
		printf("no input, return\n");
		return 0;
	}
	int ret;
	float result;
	struct calc *head;
	ret = parse_add_minus(&head, argv[1]);
	if (ret != 0) {
		printf("parse error: %d\n", ret);
		return 0;
	}
	print_calc(head);
	ret = calc_add_minus(head, &result);
	if (ret < 0) {
		printf("calculate error: %d\n", ret);
		return 0;
	}
	printf("result is %f\n", result);
}
