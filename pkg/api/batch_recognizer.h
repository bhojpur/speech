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

#ifndef VOSK_BATCH_RECOGNIZER_H
#define VOSK_BATCH_RECOGNIZER_H

#include "base/kaldi-common.h"
#include "util/common-utils.h"
#include "feat/resample.h"

#include <queue>

#include "batch_model.h"

using namespace kaldi;

class BatchRecognizer {
    public:
        BatchRecognizer(BatchModel *model, float sample_frequency);
        ~BatchRecognizer();

        void AcceptWaveform(const char *data, int len);
        int GetNumPendingChunks();
        const char *FrontResult();
        void Pop();
        void FinishStream();
        void SetNLSML(bool nlsml);

    private:

        void PushLattice(CompactLattice &clat, BaseFloat offset);

        BatchModel *model_;
        uint64_t id_;
        bool initialized_;
        bool callbacks_set_;
        bool nlsml_;
        float sample_frequency_;
        std::queue<std::string> results_;
        LinearResample *resampler_;
        kaldi::Vector<BaseFloat> buffer_;
};

#endif /* VOSK_BATCH_RECOGNIZER_H */