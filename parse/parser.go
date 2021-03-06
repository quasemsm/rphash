package parse;

import (
  "errors"
  "math"
  "reflect"
  "encoding/json"
);

var (
  fixedDecimalPoint = 18
  floatType = reflect.TypeOf(float64(0))
  weightMax = math.Abs(ToFixed(math.MaxFloat64, fixedDecimalPoint))
  weightMin = float64(0)
);

type Schema struct {
  dataType reflect.Type;
  max float64;
  min float64;
};

func NewSchema(value float64) *Schema {
  return &Schema{
    dataType: reflect.TypeOf(value),
    max: value,
    min: value,
  };
}

func (this *Schema) SetMax(floatValue float64) {
  this.max = floatValue;
};

func (this *Schema) SetMin(floatValue float64) {
  this.min = floatValue;
};

func (this *Schema) GetMax() float64 {
  return this.max;
};

func (this *Schema) GetMin() float64 {
  return this.min;
};

func Round(num float64) int {
  return int(num + math.Copysign(0.5, num));
};

func ToFixed(num float64, precision int) float64 {
  output := math.Pow(10, float64(precision));
  return float64(Round(num * output)) / output;
};

func Normalize(value float64) float64 {
  return (value - weightMin) / (weightMax - weightMin);
};

func DeNormalize(normalized float64) float64 {
	return (normalized * (weightMax - weightMin) + weightMin);
};

type Parser struct {
  schemaKeys []string;
  schema map[string]*Schema;
  label string;
};

func NewParser() *Parser {
  var schemaKeys []string;
  return &Parser{
    label: "",
    schema: nil,
    schemaKeys: schemaKeys,
  };
};

// Convert an array of bytes to a JSON struct.
func (this *Parser) BytesToJSON(bytesContents []byte) map[string]interface{} {
  var data map[string]interface{}
  if err := json.Unmarshal(bytesContents, &data); err != nil {
    panic(err);
  }
  return data;
};

func (this *Parser) JSONToBytes(jsonMap interface{}) []byte {
  bytesContents, _ := json.MarshalIndent(jsonMap, "", "  ");
  return bytesContents;
};

// Convert a json object with a schema to an array of 64 bit floats.
func (this *Parser) JSONToFloat64(jsonMap map[string]interface{}) []float64 {

  // Create an array of 64 bit floats of the same size.
  result := make([]float64, len(this.schemaKeys));

  // Iterate over the json fields and assign floating point values to each field value.
  for i := 0; i < len(this.schemaKeys); i++ {
    // Normalize the mapped value
    key := this.schemaKeys[i];
    float, _ := this.ConvertInterfaceToFloat64(jsonMap[key]);
    result[i] = Normalize(float);
  }
  return result;
};

// Convert an array of 64 bit floats to JSON according to a schema.
func (this *Parser) Float64ToJSON(floats []float64) map[string]interface{} {
  // Create an JSON object.
  jsonMap := make(map[string]interface{});

  for i := 0; i < len(this.schemaKeys); i++ {
    // DeNormalize the mapped value
    jsonMap[this.schemaKeys[i]] = DeNormalize(floats[i]);
  }
  return jsonMap;
};

// Convert a JSON table to a array of float64 arrays.
// The data should come in like this:
// {
//  "label": [{
//   "field-1": "value-1",
//   ...
//  }, {
//  "field-1": "value-1",
//  ...
//  }]
// }
func (this *Parser) JSONToFloat64Matrix(label string, dataSet map[string]interface{}) [][]float64 {
  // Assign a label to the specific schema.
  this.label = label;

  // Read the data in as an array of json objects.
  data := dataSet[label].([]interface{});
  count := len(data);

  // Allocate an array of arrays for the return.
  matrix := make([][]float64, count, count);

  // Create a schema based on an entry in the data.
  this.schema = this.CreateSchema(data);

  // Convert the json data to weighted float values.
  for i := 0; i < count; i++ {
    matrix[i] = this.JSONToFloat64(data[i].(map[string]interface{}));
  }
  return matrix;
};

// Convert a matrix of 64 bit floats to JSON according to a json schema.
// label - string associated with JSON data set schema.
// data - the array of arrays associated with the entries of data.
func (this *Parser) Float64MatrixToJSON(label string, dataSet [][]float64) map[string]interface{} {
  count := len(dataSet);

  // Create an array of JSON objects.
  data := make([]interface{}, count);

  // Create a json object to hold the array of JSON objects with the specific label.
  result := make(map[string]interface{});

  // Convert the weighted float values back to the JSON using the schema.
  for i := 0; i < count; i++ {
    data[i] = this.Float64ToJSON(dataSet[i]);
  }

  // Assign the label to the json data.
  result[this.label] = data;
  return result;
};

// Convert an unknown interface to a 64 bit floating point.
// From stackoverflow.com
func (this *Parser) ConvertInterfaceToFloat64(unk interface{}) (float64, error) {
  v := reflect.ValueOf(unk);
  v = reflect.Indirect(v);
  if !v.Type().ConvertibleTo(floatType) {
      return 0, errors.New("Cannot convert" + v.Type().String() + "to float64");
  }
  fv := v.Convert(floatType);
  return fv.Float(), nil;
};

// Create a schema based on a JSON object.
func (this *Parser) CreateSchema(data []interface{}) map[string]*Schema {
  count := len(data);

  // Set up a base schema.
  schema := make(map[string]*Schema);

  // Loop over each JSON object in the array update the schema associated schema.
  for i := 0; i < count; i++ {
    // Convert the data to a json object.
    jsonMap := data[i].(map[string]interface{});

    // Loop over its key -> value pairs.
    for key, value := range jsonMap {
      floatValue, _ := this.ConvertInterfaceToFloat64(value);
      // Has the schema not been added for the key?
      if _, ok := schema[key]; !ok {
        // Assign the key associated with the JSON field to its value type max and min.
        schema[key] = NewSchema(floatValue);

        // Assure the keys are in the proper order.
        this.schemaKeys = append(this.schemaKeys, key);
        continue;
      }

      // Check if the next value is less than the current minimum
      // Check if the next value is greater than the current maximum
      if floatValue < schema[key].GetMin() {
        schema[key].SetMin(floatValue);
      } else if floatValue > schema[key].GetMax() {
        schema[key].SetMax(floatValue);
      }
    }
  }
  return schema;
};
