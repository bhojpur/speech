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

#include "spk_model.h"

SpkModel::SpkModel(const char *speaker_path) {
    std::string speaker_path_str(speaker_path);

    ReadConfigFromFile(speaker_path_str + "/mfcc.conf", &spkvector_mfcc_opts);
    spkvector_mfcc_opts.frame_opts.allow_downsample = true; // It is safe to downsample

    ReadKaldiObject(speaker_path_str + "/final.ext.raw", &speaker_nnet);
    SetBatchnormTestMode(true, &speaker_nnet);
    SetDropoutTestMode(true, &speaker_nnet);
    CollapseModel(nnet3::CollapseModelConfig(), &speaker_nnet);

    ReadKaldiObject(speaker_path_str + "/mean.vec", &mean);
    ReadKaldiObject(speaker_path_str + "/transform.mat", &transform);

    ref_cnt_ = 1;
}

void SpkModel::Ref()
{
    std::atomic_fetch_add_explicit(&ref_cnt_, 1, std::memory_order_relaxed);
}

void SpkModel::Unref()
{
    if (std::atomic_fetch_sub_explicit(&ref_cnt_, 1, std::memory_order_release) == 1) {
         std::atomic_thread_fence(std::memory_order_acquire);
         delete this;
    }
}