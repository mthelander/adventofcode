def test(a, i)
  (i ^ 2) ^ ((a >> (i ^ 2)) ^ 3)
end

def search(out, x=0)
  return [x] if out.empty?
  x <<= 3
  z = out[-1]
  vals = []
  8.times do |i|
    a = x ^ i
    r = test(a, i)
    if r % 8 == z
      v = search(out[0..-2], a)
      unless v.empty?
        return vals + v
      end
    end
  end
  return vals
end

p search([2,4,1,2,7,5,1,3,4,3,5,5,0,3,3,0])
