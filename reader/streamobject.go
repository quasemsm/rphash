package reader;

import (
    "github.com/wenkesj/rphash/decoder"
    "github.com/wenkesj/rphash/types"
    "github.com/wenkesj/rphash/utils"
);

type StreamObject struct {
    data types.Iterator;
    numberOfProjections int;
    decoderMultiplier int;
    randomSeed int64;
    numberOfBlurs int;
    k int;
    dimension int;
    hashModulus int64;
    centroids [][]float64;
    topIDs []int64;
    decoder types.Decoder;
};

func NewStreamObject(dimension, k int) *StreamObject {
    innerDecoder := decoder.InnerDecoder();
    decoderMultiplier := 1;
    decoder := decoder.NewMultiDecoder(decoderMultiplier * innerDecoder.GetDimensionality(), innerDecoder);
    var centroids [][]float64;
    var topIDs []int64;
    return &StreamObject{
        decoder: decoder,
        dimension: dimension,
        randomSeed: int64(0),
        hashModulus: 2147483647,
        decoderMultiplier: decoderMultiplier,
        numberOfProjections: 2,
        numberOfBlurs: 2,
        k: k,
        data: nil,
        topIDs: topIDs,
        centroids: centroids,
    };
};

func (this *StreamObject) GetK() int {
    return this.k;
};
func (this *StreamObject) NumDataPoints() int {
  //JF TODO count the streamed data
  return 0;
}

func (this *StreamObject) GetDimensions() int {
    return this.dimension;
};

func (this *StreamObject) GetRandomSeed() int64 {
    return this.randomSeed;
};

func (this *StreamObject) GetNumberOfBlurs() int {
    return this.numberOfBlurs;
};

func (this *StreamObject) GetVectorIterator() types.Iterator {
    return this.data;
};

func (this *StreamObject) GetCentroids() [][]float64 {
    return this.centroids;
};

func (this *StreamObject) GetPreviousTopID() []int64 {
    return this.topIDs;
};

func (this *StreamObject) SetPreviousTopID(top []int64) {
    this.topIDs = top;
};

func (this *StreamObject) AddCentroid(v []float64) {
    this.centroids = append(this.centroids, v);
};

func (this *StreamObject) SetCentroids(l [][]float64) {
    this.centroids = l;
};

func (this *StreamObject) GetNumberOfProjections() int {
    return this.numberOfProjections;
};

func (this *StreamObject) SetNumberOfProjections(probes int) {
    this.numberOfProjections = probes;
};

func (this *StreamObject) SetNumberOfBlurs(parseInt int) {
    this.numberOfBlurs = parseInt;
};

func (this *StreamObject) SetRandomSeed(parseLong int64) {
    this.randomSeed = parseLong;
};

func (this *StreamObject) GetHashModulus() int64 {
    return this.hashModulus;
};

func (this *StreamObject) SetHashModulus(parseLong int64) {
    this.hashModulus = int64(parseLong);
};

func (this *StreamObject) SetDecoderType(dec types.Decoder) {
    this.decoder = dec;
};

func (this *StreamObject) GetDecoderType() types.Decoder {
    return this.decoder;
};

func (this *StreamObject) SetVariance(data [][]float64) {
    this.decoder.SetVariance(utils.VarianceSample(data, 0.01));
};

func (this *StreamObject) GetVariance() float64 {
    return this.decoder.GetVariance();
};
