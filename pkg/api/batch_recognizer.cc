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

#include "batch_recognizer.h"

#include "fstext/fstext-utils.h"
#include "lat/sausages.h"
#include "json.h"

BatchRecognizer::BatchRecognizer(BatchModel *model, float
                                 sample_frequency) : model_(model), sample_frequency_(sample_frequency),
                                 initialized_(false), callbacks_set_(false), nlsml_(false) {
    id_ = model->GetID(this);


    resampler_ = new LinearResample(
        sample_frequency, 16000.0f,
        std::min(sample_frequency / 2, 16000.0f / 2), 6);
}

BatchRecognizer::~BatchRecognizer() {
    delete resampler_;
    // Drop the ID
}

void BatchRecognizer::FinishStream()
{
    SubVector<BaseFloat> chunk = buffer_.Range(0, buffer_.Dim());
    model_->dynamic_batcher_->Push(id_, !initialized_, true, chunk);
}

void BatchRecognizer::PushLattice(CompactLattice &clat, BaseFloat offset)
{
    fst::ScaleLattice(fst::GraphLatticeScale(0.9), &clat);

    CompactLattice aligned_lat;
    WordAlignLattice(clat, *model_->trans_model_, *model_->winfo_, 0, &aligned_lat);

    MinimumBayesRisk mbr(aligned_lat);
    const vector<BaseFloat> &conf = mbr.GetOneBestConfidences();
    const vector<int32> &words = mbr.GetOneBest();
    const vector<pair<BaseFloat, BaseFloat> > &times =
          mbr.GetOneBestTimes();

    int size = words.size();

    if (nlsml_) {

        std::stringstream ss;
        std::stringstream text;
        ss << "<?xml version=\"1.0\"?>\n";
        ss << "<result grammar=\"default\">\n";
        BaseFloat confidence = 0.0;
        for (int i = 0; i < size; i++) {
            if (i) {
                text << " ";
            }
            confidence += conf[i];
            text << model_->word_syms_->Find(words[i]);
        }
        confidence /= size;

        ss << "<interpretation grammar=\"default\" confidence=\"" << confidence << "\">\n";
        ss << "<input mode=\"speech\">" << text.str() << "</input>\n";
        ss << "<instance>" << text.str() << "</instance>\n";
        ss << "</interpretation>\n";
        ss << "</result>\n";

        results_.push(ss.str());

    } else {
        json::JSON obj;
        stringstream text;

        // Create JSON object
        for (int i = 0; i < size; i++) {
            json::JSON word;

            word["word"] = model_->word_syms_->Find(words[i]);
            word["start"] = round(times[i].first) * 0.03 + offset;
            word["end"] = round(times[i].second) * 0.03 + offset;
            word["conf"] = conf[i];
            obj["result"].append(word);

            if (i) {
                text << " ";
            }
            text << model_->word_syms_->Find(words[i]);
        }
        obj["text"] = text.str();

//      KALDI_LOG << "Result " << id << " " << obj.dump();

        results_.push(obj.dump());
    }
}

void BatchRecognizer::SetNLSML(bool nlsml)
{
    nlsml_ = nlsml;
}


void BatchRecognizer::AcceptWaveform(const char *data, int len)
{
    uint64_t id = id_;
    if (!callbacks_set_) {
        // Define the callback for results.
#if 0
         model_->cuda_pipeline_->SetBestPathCallback(
          id,
          [&, id](const std::string &str, bool partial,
                       bool endpoint_detected) {
              if (partial) {
                  KALDI_LOG << "id #" << id << " [partial] : " << str << ":";
              }

              if (endpoint_detected) {
                  KALDI_LOG << "id #" << id << " [endpoint detected]";
              }

              if (!partial) {
                  KALDI_LOG << "id #" << id << " : " << str;
              }
            });
#endif
        model_->cuda_pipeline_->SetLatticeCallback(
          id,
          [&, id](SegmentedLatticeCallbackParams& params) {
              if (params.results.empty()) {
                  KALDI_WARN << "Empty result for callback";
                  return;
              }
              CompactLattice *clat = params.results[0].GetLatticeResult();
              BaseFloat offset = params.results[0].GetTimeOffsetSeconds();
              PushLattice(*clat, offset);
          },
          CudaPipelineResult::RESULT_TYPE_LATTICE);
        callbacks_set_ = true;
    }

    Vector<BaseFloat> input_wave(len / 2);
    for (int i = 0; i < len / 2; i++)
        input_wave(i) = *(((short *)data) + i);

    Vector<BaseFloat> resampled_wave;
    resampler_->Resample(input_wave, true, &resampled_wave);

    int32 end = buffer_.Dim();
    buffer_.Resize(end + resampled_wave.Dim(), kCopyData);
    buffer_.Range(end, resampled_wave.Dim()).CopyFromVec(resampled_wave);

    // Pick chunks and submit them to the batcher
    int32 i = 0;
    while (i + model_->samples_per_chunk_ <= buffer_.Dim()) {
        model_->dynamic_batcher_->Push(id_, !initialized_, false,
                                       buffer_.Range(i, model_->samples_per_chunk_));
        initialized_ = true;
        i += model_->samples_per_chunk_;
    }

    // Keep remaining data
    if (i > 0) {
        int32 tail = buffer_.Dim() - i;
        for (int j = 0; j < tail; j++) {
            buffer_(j) = buffer_(i + j);
        }
        buffer_.Resize(tail, kCopyData);
    }
}

const char* BatchRecognizer::FrontResult()
{
    if (results_.empty()) {
        return "";
    }
    return results_.front().c_str();
}

void BatchRecognizer::Pop()
{
    if (results_.empty()) {
        return;
    }
    results_.pop();
}

int BatchRecognizer::GetNumPendingChunks()
{
    return model_->dynamic_batcher_->GetNumPendingChunks(id_);
}