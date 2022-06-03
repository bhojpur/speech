package render

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

func (rd *Render) normalizeValue() (float32, float32) {
	max := float32(0)
	min := float32(0)
	for _, v := range rd.buffer {
		if max < v {
			max = v
		}
		if min > v {
			min = v
		}
	}
	rd.normalBuffer = make([]float32, len(rd.buffer))
	for i := 0; i < len(rd.buffer); i++ {
		rd.normalBuffer[i] = rd.buffer[int(float32(i)*rd.zoomParam)] - min
		rd.normalBuffer[i] = rd.normalBuffer[i] / (max - min) // 0~1
		//fmt.Printf("normalize: %f", rd.normalBuffer[i])
		rd.normalBuffer[i] = rd.normalBuffer[i] + 0.5 //0.5~1.5
		rd.normalBuffer[i] = rd.normalBuffer[i] * float32(rd.activeHeight/2)

		//fmt.Printf(" %f\n", rd.normalBuffer[i])
	}
	return min, max
}

func (rd *Render) addBuffer(b float32) {
	rd.buffer = append(rd.buffer, float32(b))
	if len(rd.buffer) >= int(float32(rd.activeWidth)) {
		rd.buffer = rd.buffer[1:]
	}
}
