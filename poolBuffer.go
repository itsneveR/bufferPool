package poolBuffer

/*
question :
 why do we need a bytebufferpool when we have a sync.pool ?

A bytebufferpool is a pool of byte buffers that can be reused and recycled to reduce memory allocation and garbage collection overhead.
A sync.Pool is a general-purpose pool of temporary objects that can be saved and retrieved by multiple goroutines.

The main difference between a bytebufferpool and a sync.Pool
is that a bytebufferpool has more control over the size and capacity of the byte buffers,
while a sync.Pool does not guarantee any properties of the objects returned from the pool.
A bytebufferpool can also prevent memory waste due to fragmentation by limiting the maximum total size of the byte buffers in concurrent use.


 // bytebuffer is an abstraction over sync.pool
*/
import (
	"io"
	"sync"
)

type byteSlice []byte

type Buffer struct {
	Buffer byteSlice //empty byte buffer
}

func (b *Buffer) Len() int {
	return len(b.Buffer)
}

// reads(loads) the given data to the Buffer
func (b *Buffer) ReadData(r io.Reader) (int64, error) {
	//This line assigns the byte slice from the Buffer struct
	//to the local variable "initBuffer". initBuffer will be used for reading data into.
	initBuffer := b.Buffer //initial buffer

	bufferInitLen := int64(len(initBuffer)) //buffer initial length

	bufferMaxCap := int64(cap(initBuffer)) // buffer max capacity

	//if the bytebuffer has no capacity it means it's useless, so we create a new one
	//else if it's not 0, give me the largest buffer slice possible : buffer[:bufferMaxCap]
	if bufferMaxCap == 0 {

		bufferMaxCap = 64
		initBuffer = make([]byte, bufferMaxCap)

	} else {
		initBuffer = initBuffer[:bufferMaxCap]
	}

	// now we created a buffer and we want to read from reader and write into the buffer slice.

	bufferCurrentLen := bufferInitLen //int64 .. don't confuse with var buffer

	for {
		//at the start of the loop check whether the we reached to the max capacity of slice or not
		//if yes, double the size of buffer. This ensures that the ByteBuffer can accommodate more data if needed.
		if bufferCurrentLen == bufferMaxCap {
			bufferMaxCap = bufferMaxCap * 2
			newBuffer := make([]byte, bufferMaxCap)
			copy(newBuffer, initBuffer)
			initBuffer = newBuffer
		}

		//calling the Read method of io.Reader interface, it reads the data from (r)
		//into the (buffer) starting from the current position [bufferCurrentLen:]
		s := initBuffer[bufferCurrentLen:]

		byteReadCount, err := r.Read(s)

		bufferCurrentLen += int64(byteReadCount) // updating the current buffer length

		//If there's an error, it updates the Buffer's ByteBuffer field to point to the new slice p,
		//which might have a larger capacity, and returns the total number of bytes read so far (n) and the error.
		if err != nil {

			b.ByteBuffer = initBuffer[:bufferCurrentLen]

			bufferCurrentLen -= bufferInitLen

			if err == io.EOF {
				return bufferCurrentLen, nil
			}

			return bufferCurrentLen, err
		}

	}

	/*
		In summary, this function defines method on Buffer type for getting its length
		and reading binary data from an io.Reader into the ByteBuffer,
		dynamically resizing it as needed to accommodate the incoming data.
		This is useful for building up binary buffers efficiently,
		especially when the size of the data is not known in advance.
	*/
}

/**********************************/
type Pool struct {
	calls       [steps]uint64
	calibrating uint64

	defaultSize uint64
	maxSize     uint64

	SyncPool sync.Pool
}

func (b *Buffer) WriteData(w io.Writer) (int64, error) {
	WrittenBytesLen, err := w.Write(b.ByteBuffer)
	return int64(WrittenBytesLen), err

}

func (b *Buffer) Write(p []byte) (int, error) {
	/*
		We write p... in the Write() function to unpack the slice p into its individual elements.
		This is because the append() function expects a variable number of arguments,
		and we want to be able to pass in any number of bytes to be appended to the buffer.
	*/
	b.ByteBuffer = append(b.ByteBuffer, p...)
	return len(p), nil

}
