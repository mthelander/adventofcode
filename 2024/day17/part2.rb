def test(a, i)
  (i ^ 2) ^ ((a >> (i ^ 2)) ^ 3) % 8
end

def search(out, x=0)
  return x if out.empty?
  x <<= 3
  z = out[-1]
  8.times do |i|
    a = x ^ i
    if z == test(a, i)
      v = search(out[0..-2], a)
      return v unless v.nil?
    end
  end
  return nil
end

p search([2,4,1,2,7,5,1,3,4,3,5,5,0,3,3,0])
