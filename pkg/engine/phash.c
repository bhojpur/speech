// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

#include "phash.h"
#include <stdio.h>
#include <assert.h>

long long hash(int n1, int n2, double *data)
{
    int i, j, k, l;
    long long x;

    assert (n1 == 32 && n2 == 64);

    double sum = 0, sum1;
    for (i = 0; i < n1; i++) {
        for (j = 0; j < n2; j++) {
           sum = sum + data[i * n2 + j];
        }
    }
    sum /= 64;
    x = 0;
    for (i = 0; i < n1; i+=4) {
        for (j = 0; j < n2; j+=8) {
            sum1 = 0;
            for (k = 0; k < 4; k++) {
                for (l = 0; l < 8; l++) {
                    sum1 = sum1 + data[(i + k) * n2 + j + l];
                }
            }
            x <<= 1;
            if (sum1 > sum) {
                x |= 0x1;
            }
        }
    }
    return x;
}